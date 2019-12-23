package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

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
		path := vars["path"]

		filename := fmt.Sprintf("%v-%v-%v", name, version, strings.ReplaceAll(path, "/", "-"))
		filename = strings.TrimSuffix(filename, "-")
		filename += ".zip"

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%v", filename))

		err := client.Archive(name, version, path, w)
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
