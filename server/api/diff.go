package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/g-harel/rejstry/internal/diff"
	"github.com/g-harel/rejstry/internal/registry"
	"github.com/g-harel/rejstry/internal/tarball"
)

type v1DiffTarDir struct {
	version string
	dir     string
	err     error
}

func V1Diff(w http.ResponseWriter, r *http.Request) {
	// Only handle requests with POST method and correct content type.
	if r.Method != http.MethodPost || r.Header.Get("Content-Type") != "application/json" {
		http.NotFound(w, r)
		return
	}

	// Parse request object.
	req := &struct {
		Registry string `json:"registry"`
		Package  string `json:"package"`
		Version  string `json:"version"`
		Compare  string `json:"compare"`
	}{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil || req.Registry == "" || req.Package == "" || req.Version == "" || req.Compare == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if req.Version == req.Compare {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	dirChan := make(chan v1DiffTarDir)
	for _, version := range []string{req.Version, req.Compare} {
		go func(v string) {
			// Create temporary working directory.
			dir, err := ioutil.TempDir("", "")
			if err != nil {
				dirChan <- v1DiffTarDir{v, "", fmt.Errorf("create temp dir: %v", err)}
				return
			}

			// Fetch package contents for given version.
			pkg, err := registry.PackageContents(req.Registry, req.Package, v)
			if err != nil {
				dirChan <- v1DiffTarDir{v, "", fmt.Errorf("fetch package: %v", err)}
				return
			}
			defer pkg.Close()

			// Write package contents to directory.
			tarball.Extract(pkg, tarball.Downloader(dir))

			dirChan <- v1DiffTarDir{v, dir, nil}
		}(version)
	}

	dirs := map[string]string{}
	for i := 0; i < 2; i++ {
		dir := <-dirChan
		if dir.err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Printf("ERROR download package '%v': %v", dir.version, err)
			return
		}
		dirs[dir.version] = dir.dir
	}

	out, err := diff.Compare(path.Join(dirs[req.Version], "package"), path.Join(dirs[req.Compare], "package"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR compare package contents: %v", err)
		return
	}

	_, err = io.Copy(w, out)
	if err != nil {
		log.Printf("ERROR copy output: %v", err)
	}

	for _, dir := range dirs {
		_ = os.RemoveAll(dir)
	}
}
