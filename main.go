package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/g-harel/rejstry/templates"
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

func redirect(pre, post string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, pre+r.URL.Path+post, http.StatusFound)
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

	// Logs are handled by the runtime in production.
	if os.Getenv("ENV") == "production" {
		log.SetOutput(ioutil.Discard)
	}

	// Show package versions.
	r.Handle("/package/{name}", redirect("", "/"))
	r.Handle("/package/{name}/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		templates.Versions(w, r, name, "")
	}))

	// Show package contents.
	r.Handle("/package/{name}/{version}", redirect("", "/"))
	r.PathPrefix("/package/{name}/{version}/").Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		version := vars["version"]
		path := strings.Join(strings.Split(r.URL.Path, "/")[4:], "/")

		// Show a directory if the path ends with a path delimiter.
		if path == "" || path[len(path)-1] == '/' {
			templates.Directory(w, r, name, version, path)
		} else {
			templates.File(w, r, name, version, path)
		}
	}))

	// Pick second version to compare to.
	r.Handle("/compare/{name}/{a}", redirect("", "/"))
	r.Handle("/compare/{name}/{a}/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		versionA := vars["a"]

		templates.Versions(w, r, name, versionA)
	}))

	// Compare package versions.
	r.Handle("/compare/{name}/{a}/{b}", redirect("", "/"))
	r.Handle("/compare/{name}/{a}/{b}/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		versionA := vars["a"]
		versionB := vars["b"]

		templates.Compare(w, r, name, versionA, versionB)
	}))

	// Static assets.
	assets := http.FileServer(http.Dir("assets"))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", assets))
	r.Handle("/favicon.ico", redirect("/assets", ""))

	// Attempt to match single path as package name.
	// Handlers registered before this point have a higher matching priority.
	r.Handle("/{package}", redirect("/package", "/"))

	// Take port number from environment if provided.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("accepting connections at :%v", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
