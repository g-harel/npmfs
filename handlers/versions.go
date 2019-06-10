package handlers

import (
	"log"
	"net/http"
	"sort"

	"github.com/g-harel/npmfs/internal/registry"
	"github.com/g-harel/npmfs/internal/util"
	"github.com/g-harel/npmfs/templates"
	"github.com/gorilla/mux"
)

// Versions handler displays all available package versions.
func Versions(client registry.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		disabled := vars["disabled"]

		// Fetch and sort version list.
		versions, latest, err := client.Versions(name)
		if err == registry.ErrNotFound {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Printf("ERROR fetch package versions: %v", err)
			return
		}
		sort.Sort(util.SemverSort(versions))

		// Render page template.
		templates.PageVersions(name, latest, disabled, versions).Handler(w, r)
	}
}
