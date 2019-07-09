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

// Redirect responds with a temporary redirect after adding the pre and postfix.
// Path params (ex. "{name}") in pre/post are looked up in "mux.Vars()".
func redirect(pre, post string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		param := regexp.MustCompile("\\{.+\\}")
		pre = param.ReplaceAllStringFunc(pre, func(match string) string {
			return vars[match[1:len(match)-1]]
		})
		post = param.ReplaceAllStringFunc(post, func(match string) string {
			return vars[match[1:len(match)-1]]
		})
		http.Redirect(w, r, pre+r.URL.Path+post, http.StatusFound)
	}
}

// Routes returns an http handler with all the routes/handlers attached.
func routes(client registry.Client) http.Handler {
	r := mux.NewRouter()

	// Add gzip middleware to all handlers.
	r.Use(gziphandler.GzipHandler)

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
	r.HandleFunc("/package/"+nameP+"", redirect("", "/"))
	r.HandleFunc("/package/"+nameP+"/", handlers.Versions(client))

	// Show package contents.
	r.HandleFunc("/package/"+nameP+"/{version}", redirect("", "/"))
	r.HandleFunc("/package/"+nameP+"/v/{version}", redirect("", "/../../{version}/"))
	r.HandleFunc("/package/"+nameP+"/v/{version}/", redirect("", "/../../{version}/"))
	r.PathPrefix("/package/" + nameP + "/{version}/" + dirP).HandlerFunc(handlers.Directory(client))
	r.PathPrefix("/package/" + nameP + "/{version}/" + fileP).HandlerFunc(handlers.File(client))

	// Pick second version to compare to.
	r.HandleFunc("/compare/"+nameP+"/{disabled}", redirect("", "/"))
	r.HandleFunc("/compare/"+nameP+"/{disabled}/", handlers.Versions(client))

	// Compare package versions.
	r.HandleFunc("/compare/"+nameP+"/{a}/{b}", redirect("", "/"))
	r.HandleFunc("/compare/"+nameP+"/{a}/{b}/", handlers.Compare(client))

	// Static assets.
	assets := http.FileServer(http.Dir("assets"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", assets))
	r.HandleFunc("/favicon.ico", redirect("/assets", ""))
	r.HandleFunc("/robots.txt", redirect("/assets", ""))

	// Attempt to match single path as package name.
	// Handlers registered before this point have a higher matching priority.
	r.HandleFunc("/"+nameP, redirect("/package", "/"))
	r.HandleFunc("/"+nameP+"/", redirect("/package", ""))

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
