package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/g-harel/rejstry/internal/middleware"
	"github.com/g-harel/rejstry/internal/registry"
	"github.com/g-harel/rejstry/internal/semver"
	"github.com/g-harel/rejstry/server/api"
)

func pageVersions(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if parts[2] == "" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(
		"templates/layout.html",
		"templates/pages/versions.html",
	)
	if err != nil {
		// TODO
		panic(err)
	}

	versions, latest, err := registry.PackageVersions("registry.npmjs.com", parts[2])
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR fetch package versions: %v", err)
		return
	}
	sort.Sort(semver.Sort(versions))

	context := &struct {
		Package  string
		Versions []string
		Latest   string
	}{
		Package:  parts[2],
		Versions: versions,
		Latest:   latest,
	}

	err = tmpl.Execute(w, context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// API paths.
	http.HandleFunc("/api/v1/diff", middleware.Log(api.V1Diff))
	http.HandleFunc("/api/v1/file", middleware.Log(api.V1File))
	http.HandleFunc("/api/v1/files", middleware.Log(api.V1Files))
	http.HandleFunc("/api/v1/versions", middleware.Log(api.V1Versions))

	// Rendered templates.
	http.HandleFunc("/package/", middleware.Log(pageVersions))

	// Static files.
	static := http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/static/", middleware.Log(static.ServeHTTP))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("accepting connections at :%v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
