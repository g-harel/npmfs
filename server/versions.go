package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/g-harel/rejstry/internal"
)

func versionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 page not found")
		return
	}

	request := &struct {
		Registry string `json:"registry"`
		Package  string `json:"package"`
	}{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	versions, latest, err := internal.PackageVersions(request.Registry, request.Package)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// TODO logging
		log.Print(err)
		return
	}
	sort.Sort(internal.SemverSort(versions))

	response := &struct {
		Latest   string   `json:"latest"`
		Versions []string `json:"versions"`
	}{
		Latest:   latest,
		Versions: versions,
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// TODO logging
		log.Print(err)
		return
	}
}
