package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/g-harel/npmfs/internal/diff"
	"github.com/g-harel/npmfs/internal/registry"
	"github.com/g-harel/npmfs/internal/tarball"
	"github.com/g-harel/npmfs/templates"
	"github.com/gorilla/mux"
)

// Compare handler displays a diff between two package versions.
func Compare(ry registry.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		versionA := vars["a"]
		versionB := vars["b"]

		if versionA == versionB {
			http.NotFound(w, r)
			return
		}

		// Download both package version contents to a temporary directory in parallel.
		type downloadedDir struct {
			version string
			dir     string
			err     error
		}
		dirChan := make(chan downloadedDir)
		for _, version := range []string{versionA, versionB} {
			go func(v string) {
				// Create temporary working directory.
				dir, err := ioutil.TempDir("", "")
				if err != nil {
					dirChan <- downloadedDir{v, "", fmt.Errorf("create temp dir: %v", err)}
					return
				}

				// Fetch package contents for given version.
				pkg, err := ry.PackageContents(name, v)
				if err != nil {
					// Error not wrapped so it can be checked against "registry.ErrNotFound".
					dirChan <- downloadedDir{v, "", err}
					return
				}
				defer pkg.Close()

				// Write package contents to directory.
				err = tarball.Extract(pkg, tarball.Downloader(func(name string) string {
					return path.Join(dir, strings.TrimPrefix(name, "package"))
				}))
				if err != nil {
					dirChan <- downloadedDir{v, "", fmt.Errorf("download contents: %v", err)}
					return
				}

				dirChan <- downloadedDir{v, dir, nil}
			}(version)
		}

		// Wait for both version's contents to be downloaded.
		dirs := map[string]string{}
		for i := 0; i < 2; i++ {
			dir := <-dirChan
			if dir.err == registry.ErrNotFound {
				http.NotFound(w, r)
				return
			}
			if dir.err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Printf("ERROR download package '%v': %v", dir.version, dir.err)
				return
			}
			dirs[dir.version] = dir.dir
		}

		// Compare contents.
		patches, err := diff.Compare(dirs[versionA], dirs[versionB])
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
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
