package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/g-harel/rejstry/internal/registry"
	"github.com/g-harel/rejstry/internal/tarball"
)

func v1Files(w http.ResponseWriter, r *http.Request) {
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
	}{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil || req.Registry == "" || req.Package == "" || req.Version == "" {
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

	// Extract files from package contents.
	files := []string{}
	err = tarball.Extract(pkg, func(name string, contents io.Reader) error {
		files = append(files, strings.TrimPrefix(name, "package/"))
		return nil
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR extract files from package contents: %v", err)
		return
	}

	// Write data to response.
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(files)
	if err != nil {
		log.Printf("ERROR encode response: %v", err)
		return
	}
}
