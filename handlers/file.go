package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/g-harel/npmfs/internal/registry"
	"github.com/g-harel/npmfs/internal/util"
	"github.com/g-harel/npmfs/templates"
	"github.com/gorilla/mux"
)

// File handler displays a file view of package contents at the provided path.
func File(client registry.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		version := vars["version"]
		path := vars["path"]

		// Fetch file contents.
		file, err := client.File(name, version, path)
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
			log.Printf("ERROR fetch file: %v", err)
			return
		}

		parts, links := util.BreakPathRelative(path)
		lines := strings.Split(file, "\n")

		// Render page template.
		templates.PageFile(name, version, parts, links, lines).Handler(w, r)
	}
}
