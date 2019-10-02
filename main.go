package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/NYTimes/gziphandler"
	"github.com/g-harel/npmfs/handlers"
	"github.com/g-harel/npmfs/internal/registry"
	"github.com/g-harel/npmfs/internal/registry/standard"
	"github.com/g-harel/npmfs/templates"
	"github.com/gorilla/mux"
)

// Known public registries.
var (
	NPM  registry.Client = &standard.Client{Host: "registry.npmjs.com"}
	Yarn registry.Client = &standard.Client{Host: "registry.yarnpkg.com"}
	Open registry.Client = &standard.Client{Host: "npm.open-registry.dev"}
)

// Redirect responds with a temporary redirect to the rewritten path.
// Path params (ex. "{name}") are looked up in "mux.Vars()".
func redirect(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		param := regexp.MustCompile("\\{.+?\\}")
		newPath := param.ReplaceAllStringFunc(path, func(match string) string {
			return vars[match[1:len(match)-1]]
		})
		http.Redirect(w, r, newPath, http.StatusFound)
	}
}

// HTTPSHandler is middleware to redirect HTTP requests to the HTTPS equivalent.
// The middleware assumes it is receiving a forwarded request.
// Local development is not impacted unless requests specify the checked header.
func httpsHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Forwarded-Proto") == "http" {
			newURL := fmt.Sprintf("https://%v%v", r.Host, r.RequestURI)
			http.Redirect(w, r, newURL, http.StatusMovedPermanently)
		}
		h.ServeHTTP(w, r)
	})
}

// Routes returns an http handler with all the routes/handlers attached.
func routes(client registry.Client) http.Handler {
	r := mux.NewRouter()

	// Add gzip middleware to all handlers.
	r.Use(gziphandler.GzipHandler)

	// Redirect all HTTP requests.
	r.Use(httpsHandler)

	// Show homepage.
	r.HandleFunc("/", templates.PageHome().Handler)

	var (
		// Name pattern matches with simple and org-scoped names.
		// (ex. "lodash", "react", "@types/express")
		nameP = "{name:(?:@[^/]+\\/)?[^/]+}"
		// Directory path pattern matches everything that ends with a path separator.
		dirP = "{path:(?:.+/)?$}"
		// File path pattern matches everything that does not end in a path separator.
		fileP = "{path:.*[^/]$}"
	)

	// Show package versions.
	r.HandleFunc("/package/"+nameP+"/", handlers.Versions(client))
	r.HandleFunc("/package/"+nameP, redirect("/package/{name}/"))
	r.HandleFunc("/compare/"+nameP+"/", redirect("/package/{name}/"))
	r.HandleFunc("/compare/"+nameP, redirect("/package/{name}/"))

	// Show package contents.
	r.PathPrefix("/package/" + nameP + "/{version}/" + dirP).HandlerFunc(handlers.Directory(client))
	r.PathPrefix("/package/" + nameP + "/{version}/" + fileP).HandlerFunc(handlers.File(client))
	r.HandleFunc("/package/"+nameP+"/{version}", redirect("/package/{name}/{version}/"))
	r.HandleFunc("/package/"+nameP+"/v/{version}", redirect("/package/{name}/{version}/"))
	r.HandleFunc("/package/"+nameP+"/v/{version}/", redirect("/package/{name}/{version}/"))

	// Pick second version to compare to.
	r.HandleFunc("/compare/"+nameP+"/{disabled}/", handlers.Versions(client))
	r.HandleFunc("/compare/"+nameP+"/{disabled}", redirect("/compare/{name}/{disabled}/"))

	// Compare package versions.
	r.HandleFunc("/compare/"+nameP+"/{a}/{b}/", handlers.Compare(client))
	r.HandleFunc("/compare/"+nameP+"/{a}/{b}", redirect("/compare/{name}/{a}/{b}/"))

	// Download package contents.
	r.HandleFunc("/download/"+nameP+"/{version}/", handlers.Download(client))
	r.HandleFunc("/download/"+nameP+"/{version}", redirect("/download/{name}/{version}/"))

	// Static assets.
	assets := http.FileServer(http.Dir("assets"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", assets))
	r.HandleFunc("/favicon.ico", redirect("/assets/favicon.ico"))
	r.HandleFunc("/robots.txt", redirect("/assets/robots.txt"))

	// Attempt to match single path as package name.
	// Handlers registered before this point have a higher matching priority.
	r.HandleFunc("/"+nameP, redirect("/package/{name}/"))
	r.HandleFunc("/"+nameP+"/", redirect("/package/{name}/"))

	return r
}

func main() {
	// Take port number from environment if provided.
	// https://cloud.google.com/run/docs/reference/container-contract
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("accepting connections at :%v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), routes(NPM)))
}
