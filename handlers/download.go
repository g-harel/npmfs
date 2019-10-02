package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/g-harel/npmfs/internal/registry"
	"github.com/g-harel/npmfs/templates"
	"github.com/gorilla/mux"
	"golang.org/x/xerrors"
)

// Download handler serves a zip archive of the package contents.
func Download(client registry.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		version := vars["version"]

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%v-%v.zip", name, version))

		err := client.Archive(name, version, w)
		if err != nil {
			var registryErr *registry.Error
			if xerrors.As(err, &registryErr) {
				templates.PageError(registryErr.StatusCode, registryErr.Error()).Handler(w, r)
				return
			}
			templates.PageError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)).Handler(w, r)
			log.Printf("ERROR fetch package archive: %v", err)
			return
		}
	}
}
