package handlers

import (
	"log"
	"net/http"
	"sort"

	"github.com/g-harel/npmfs/internal/registry"
	"github.com/g-harel/npmfs/internal/util"
	"github.com/g-harel/npmfs/templates"
	"github.com/gorilla/mux"
	"golang.org/x/xerrors"
)

// Versions handler displays all available package versions.
func Versions(client registry.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		disabled := vars["disabled"]

		// Fetch and sort version list.
		versions, latest, err := client.Versions(name)
		if err != nil {
			var registryErr *registry.Err
			if xerrors.As(err, &registryErr) {
				http.Error(w, registryErr.Error(), registryErr.StatusCode)
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Printf("ERROR fetch package versions: %v", err)
			return
		}
		sort.Sort(util.SemverSort(versions))

		// Render page template.
		templates.PageVersions(name, latest, disabled, versions).Handler(w, r)
	}
}
