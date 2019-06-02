package handlers

import (
	"log"
	"net/http"
	"sort"

	"github.com/g-harel/npmfs/internal/registry"
	"github.com/g-harel/npmfs/internal/semver"
	"github.com/g-harel/npmfs/templates"
	"github.com/gorilla/mux"
)

func Versions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	disabled := vars["disabled"]

	versions, latest, err := registry.PackageVersions("registry.npmjs.com", name)
	if err == registry.ErrNotFound {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR fetch package versions: %v", err)
		return
	}
	sort.Sort(semver.Sort(versions))

	templates.PageVersions(name, latest, disabled, versions).Render(w)
}
