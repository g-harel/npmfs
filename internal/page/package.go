package page

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/g-harel/rejstry/internal/registry"
	"github.com/g-harel/rejstry/internal/semver"
	"github.com/g-harel/rejstry/internal/tarball"
)

func Package() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// First two parts are removed because they are never used (/package/).
		parts := strings.SplitN(r.URL.Path, "/", 5)[2:]

		// Package name must be provided.
		if len(parts) < 1 {
			http.NotFound(w, r)
			return
		}

		name := parts[0]

		// Package name must be provided.
		if name == "" {
			http.NotFound(w, r)
			return
		}

		// Enforce trailing slash after package name (for relative links).
		if len(parts) < 2 {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusFound)
			return
		}

		version := parts[1]

		// Handle requests that don't specify a version.
		if version == "" {
			packageVersions(name)(w, r)
			return
		}

		// Enforce trailing slash after package version (for relative links).
		if len(parts) < 3 {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusFound)
			return
		}

		file := parts[2]

		// Handle requests that don't specify a file.
		if file == "" {
			packageFiles(name, version)(w, r)
			return
		}

		// TODO
		log.Printf("=== LIST FILE (%v) ===", file)
		http.NotFound(w, r)
	})
}

func packageVersions(name string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(
			"templates/layout.html",
			"templates/pages/versions.html",
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("ERROR parse template: %v", err)
			return
		}

		versions, latest, err := registry.PackageVersions("registry.npmjs.com", name)
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
			Package:  name,
			Versions: versions,
			Latest:   latest,
		}

		err = tmpl.Execute(w, context)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("ERROR execute template: %v", err)
			return
		}
	})
}

func packageFiles(name, version string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(
			"templates/layout.html",
			"templates/pages/files.html",
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("ERROR parse template: %v", err)
			return
		}

		// Fetch package contents.
		pkg, err := registry.PackageContents("registry.npmjs.com", name, version)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Printf("ERROR fetch package contents: %v", err)
			return
		}
		defer pkg.Close()

		// Extract files from package contents.
		files := []string{}
		err = tarball.Extract(pkg, func(name string, contents io.Reader) error {
			files = append(files, strings.TrimPrefix(name, "package/"))
			return nil
		})
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Printf("ERROR extract files from package contents: %v", err)
			return
		}

		context := &struct {
			Package string
			Version string
			Files   []string
		}{
			Package: name,
			Version: version,
			Files:   files,
		}

		err = tmpl.Execute(w, context)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("ERROR execute template: %v", err)
			return
		}
	})
}
