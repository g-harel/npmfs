package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/g-harel/npmfs/internal/diff"
	"github.com/g-harel/npmfs/internal/registry"
	"github.com/g-harel/npmfs/templates"
	"github.com/gorilla/mux"
	"golang.org/x/xerrors"
)

// Compare handler displays a diff between two package versions.
func Compare(client registry.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		versionA := vars["a"]
		versionB := vars["b"]

		if versionA == versionB {
			templates.PageError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest)).Handler(w, r)
			return
		}
		versions := []string{versionA, versionB}

		// Download both package version contents in parallel.
		type downloadedDir struct {
			version string
			dir     string
			err     error
		}
		dirChan := make(chan downloadedDir)
		for _, version := range versions {
			go func(v string) {
				dir, err := client.Download(name, v)
				dirChan <- downloadedDir{v, dir, err}
			}(version)
		}

		// Wait for both version's contents to be downloaded.
		dirs := map[string]string{}
		for _ = range versions {
			dir := <-dirChan
			if dir.err != nil {
				var registryErr *registry.Error
				if xerrors.As(dir.err, &registryErr) {
					templates.PageError(registryErr.StatusCode, registryErr.Error()).Handler(w, r)
					return
				}
				templates.PageError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)).Handler(w, r)
				log.Printf("ERROR download package '%v': %v", dir.version, dir.err)
				return
			}
			dirs[dir.version] = dir.dir
		}

		// Compare contents.
		patches, err := diff.Compare(dirs[versionA], dirs[versionB])
		if err != nil {
			templates.PageError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)).Handler(w, r)
			log.Printf("ERROR compare package contents: %v", err)
			return
		}

		// Cleanup created directories.
		for _, path := range dirs {
			_ = os.RemoveAll(path)
		}

		// Render page template.
		templates.PageCompare(name, versionA, versionB, patches).Handler(w, r)
	}
}
