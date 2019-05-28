package page

import (
	"html/template"
	"log"
	"net/http"
	"sort"

	"github.com/g-harel/rejstry/internal/registry"
	"github.com/g-harel/rejstry/internal/semver"
)

func Versions(w http.ResponseWriter, r *http.Request, name string) {
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
}
