package handlers

import (
	"log"
	"net/http"
	"strings"

	"github.com/g-harel/npmfs/internal/registry"
	"github.com/g-harel/npmfs/internal/util"
	"github.com/g-harel/npmfs/templates"
	"github.com/gorilla/mux"
	"golang.org/x/xerrors"
)

// ReadFile handler displays a file view of package contents at the provided path.
func ReadFile(client registry.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		version := vars["version"]
		path := vars["path"]

		// Fetch file contents.
		file, err := client.File(name, version, path)
		if err != nil {
			var registryErr *registry.Error
			if xerrors.As(err, &registryErr) {
				templates.PageError(registryErr.StatusCode, registryErr.Error()).Handler(w, r)
				return
			}
			templates.PageError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)).Handler(w, r)
			log.Printf("ERROR fetch file: %v", err)
			return
		}

		parts, links := util.BreakPathRelative(path)
		lines := strings.Split(file, "\n")

		// Render page template.
		templates.PageFile(name, version, util.ByteCount(int64(len(file))), parts, links, lines).Handler(w, r)
	}
}

// ReadDir handler displays a directory view of package contents at the provided path.
func ReadDir(client registry.Client) http.HandlerFunc {
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
