package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/g-harel/rejstry/internal/registry"
	"github.com/g-harel/rejstry/internal/tarball"
)

func V1File(w http.ResponseWriter, r *http.Request) {
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
		Path     string `json:"path"`
	}{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil || req.Registry == "" || req.Package == "" || req.Version == "" || req.Path == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Fetch package contents.
	pkg, err := registry.PackageContents(req.Registry, req.Package, req.Version)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR fetch package contents: %v", err)
		return
	}
	defer pkg.Close()

	// Write file contents to response.
	found := false
	err = tarball.Extract(pkg, func(name string, contents io.Reader) error {
		if strings.TrimPrefix(name, "package/") == req.Path {
			found = true
			_, err := io.Copy(w, contents)
			if err != nil {
				log.Printf("ERROR copy contents: %v", err)
			}
		}
		return nil
	})
	if err != nil {
		log.Printf("ERROR extract files from package contents: %v", err)
		return
	}
	if !found {
		http.NotFound(w, r)
	}
}
