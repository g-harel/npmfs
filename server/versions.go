package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"github.com/g-harel/rejstry/internal"
)

func v1Versions(w http.ResponseWriter, r *http.Request) {
	// Only handle requests with POST method and correct content type.
	if r.Method != http.MethodPost || r.Header.Get("Content-Type") != "application/json" {
		http.NotFound(w, r)
		return
	}

	// Parse request object.
	req := &struct {
		Registry string `json:"registry"`
		Package  string `json:"package"`
	}{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil || req.Registry == "" || req.Package == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Fetch version data.
	versions, latest, err := internal.PackageVersions(req.Registry, req.Package)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR fetch package versions: %v", err)
		return
	}
	sort.Sort(internal.SemverSort(versions))

	// Write data to response.
	res := &struct {
		Latest   string   `json:"latest"`
		Versions []string `json:"versions"`
	}{
		Latest:   latest,
		Versions: versions,
	}
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Printf("ERROR encode response: %v", err)
		return
	}
}
