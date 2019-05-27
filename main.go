package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/g-harel/rejstry/internal/page"
	"github.com/gorilla/mux"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(statusCode int) {
	s.status = statusCode
	s.ResponseWriter.WriteHeader(statusCode)
}

func redirect(append string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.URL.Path+"/", http.StatusMovedPermanently)
	})
}

func main() {
	r := mux.NewRouter()

	// Add logging middleware.
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{w, http.StatusOK}
			h.ServeHTTP(rec, r)
			log.Printf("%v %v %v - %vms", r.Method, r.RequestURI, rec.status, int64(time.Since(start)/time.Millisecond))
		})
	})

	r.Handle("/package/{name}", redirect("/"))
	r.Handle("/package/{name}/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		page.PackageVersions(w, r, name)
	}))

	r.Handle("/package/{name}/{version}", redirect("/"))
	r.Handle("/package/{name}/{version}/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		version := vars["version"]

		page.PackageFiles(w, r, name, version)
	}))

	// Static files.
	static := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", static))

	// Take port number from environment if provided.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("accepting connections at :%v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
