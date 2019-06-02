package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"github.com/g-harel/npmfs/internal/diff"
	"github.com/g-harel/npmfs/internal/registry"
	"github.com/g-harel/npmfs/internal/tarball"
	"github.com/g-harel/npmfs/templates"
	"github.com/gorilla/mux"
)

func Compare(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	versionA := vars["a"]
	versionB := vars["b"]

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
			pkg, err := registry.PackageContents("registry.npmjs.com", name, v)
			if err != nil {
				// Error not wrapped so it can be checked against "registry.ErrNotFound".
				dirChan <- downloadedDir{v, "", err}
				return
			}
			defer pkg.Close()

			// Write package contents to directory.
			tarball.Extract(pkg, tarball.Downloader(dir))

			dirChan <- downloadedDir{v, dir, nil}
		}(version)
	}

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

	patches, err := diff.Compare(path.Join(dirs[versionA], "package"), path.Join(dirs[versionB], "package"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR compare package contents: %v", err)
		return
	}

	templates.PageCompare(name, versionA, versionB, patches).Render(w)
}
