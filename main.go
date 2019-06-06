package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/NYTimes/gziphandler"
	"github.com/g-harel/npmfs/handlers"
	"github.com/g-harel/npmfs/templates"
	"github.com/gorilla/mux"
)

// Name pattern matches with simple and org-scoped names.
// (ex. "lodash", "react", "@types/express")
const namePattern = "{name:(?:@[^/]+\\/)?[^/]+}"

// Directory path pattern matches everything that ends with a path separator.
const dirPattern = "{path:(?:.+/)?$}"

// File path pattern matches everything that does not end in a path separator.
const filePattern = "{path:.*[^/]$}"

// Redirect responds with a temporary redirect after adding the pre and postfix.
func redirect(pre, post string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, pre+r.URL.Path+post, http.StatusFound)
	}
}

func main() {
	r := mux.NewRouter()

	// Add gzip middleware to all handlers.
	r.Use(gziphandler.GzipHandler)

	// Show homepage.
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templates.PageHome().Render(w)
	})

	// Show package versions.
	r.HandleFunc("/package/"+namePattern+"", redirect("", "/"))
	r.HandleFunc("/package/"+namePattern+"/", handlers.Versions)

	// Show package contents.
	r.HandleFunc("/package/"+namePattern+"/{version}", redirect("", "/"))
	r.PathPrefix("/package/" + namePattern + "/{version}/" + dirPattern).HandlerFunc(handlers.Directory)
	r.PathPrefix("/package/" + namePattern + "/{version}/" + filePattern).HandlerFunc(handlers.File)

	// Pick second version to compare to.
	r.HandleFunc("/compare/"+namePattern+"/{disabled}", redirect("", "/"))
	r.HandleFunc("/compare/"+namePattern+"/{disabled}/", handlers.Versions)

	// Compare package versions.
	r.HandleFunc("/compare/"+namePattern+"/{a}/{b}", redirect("", "/"))
	r.HandleFunc("/compare/"+namePattern+"/{a}/{b}/", handlers.Compare)

	// Static assets.
	assets := http.FileServer(http.Dir("assets"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", assets))
	r.HandleFunc("/favicon.ico", redirect("/assets", ""))
	r.HandleFunc("/robots.txt", redirect("/assets", ""))

	// Attempt to match single path as package name.
	// Handlers registered before this point have a higher matching priority.
	r.HandleFunc("/"+namePattern, redirect("/package", "/"))
	r.HandleFunc("/"+namePattern+"/", redirect("/package", ""))

	// Take port number from environment if provided.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("accepting connections at :%v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), r))
}
