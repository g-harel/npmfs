package handlers

import (
	"log"
	"net/http"

	"github.com/g-harel/npmfs/internal/registry"
	"github.com/g-harel/npmfs/internal/util"
	"github.com/g-harel/npmfs/templates"
	"github.com/gorilla/mux"
	"golang.org/x/xerrors"
)

// Directory handler displays a directory view of package contents at the provided path.
func Directory(client registry.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		version := vars["version"]
		path := vars["path"]

		dirs, files, err := client.Directory(name, version, path)
		if err != nil {
			var registryErr *registry.Error
			if xerrors.As(err, &registryErr) {
				templates.PageError(registryErr.StatusCode, registryErr.Error()).Handler(w, r)
				return
			}
			templates.PageError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)).Handler(w, r)
			log.Printf("ERROR fetch directory: %v", err)
			return
		}

		parts, links := util.BreakPathRelative(path)

		// Render page template.
		templates.PageDirectory(name, version, parts, links, dirs, files).Handler(w, r)
	}
}
