package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/g-harel/rejstry/internal/middleware"
	"github.com/g-harel/rejstry/server/api"
)

func pageVersions(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/layout.html",
		"templates/pages/versions.html",
	)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(w, nil)
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
	http.HandleFunc("/versions", middleware.Log(pageVersions))

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
