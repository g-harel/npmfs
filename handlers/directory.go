package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/g-harel/npmfs/internal/registry"
	"github.com/g-harel/npmfs/internal/util"
	"github.com/g-harel/npmfs/templates"
	"github.com/gorilla/mux"
)

// Directory handler displays a directory view of package contents at the provided path.
func Directory(client registry.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		version := vars["version"]
		path := vars["path"]

		dirs, files, err := client.Directory(name, version, path)
		if err == registry.ErrNotFound {
			http.NotFound(w, r)
			return
		}
		if err == registry.ErrTimeout {
			http.Error(w, fmt.Sprintf("%v timeout", http.StatusGatewayTimeout), http.StatusGatewayTimeout)
			return
		}
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Printf("ERROR fetch directory: %v", err)
			return
		}

		parts, links := util.BreakPathRelative(path)

		// Render page template.
		templates.PageDirectory(name, version, parts, links, dirs, files).Handler(w, r)
	}
}
