package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/g-harel/rejstry/handlers"
	"github.com/gorilla/mux"
)

func redirect(pre, post string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, pre+r.URL.Path+post, http.StatusFound)
	}
}

func main() {
	r := mux.NewRouter()

	// Show package versions.
	r.HandleFunc("/package/{name}", redirect("", "/"))
	r.HandleFunc("/package/{name}/", handlers.Versions)

	// Show package contents.
	r.HandleFunc("/package/{name}/{version}", redirect("", "/"))
	r.PathPrefix("/package/{name}/{version}/{path:(?:.+/)?$}").HandlerFunc(handlers.Directory)
	r.PathPrefix("/package/{name}/{version}/{path:.*}").HandlerFunc(handlers.File)

	// Pick second version to compare to.
	r.HandleFunc("/compare/{name}/{disabled}", redirect("", "/"))
	r.HandleFunc("/compare/{name}/{disabled}/", handlers.Versions)

	// Compare package versions.
	r.HandleFunc("/compare/{name}/{a}/{b}", redirect("", "/"))
	r.HandleFunc("/compare/{name}/{a}/{b}/", handlers.Compare)

	// Static assets.
	assets := http.FileServer(http.Dir("assets"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", assets))
	r.HandleFunc("/favicon.ico", redirect("/assets", ""))

	// Attempt to match single path as package name.
	// Handlers registered before this point have a higher matching priority.
	r.HandleFunc("/{package}", redirect("/package", "/"))

	// Take port number from environment if provided.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("accepting connections at :%v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), r))
}
