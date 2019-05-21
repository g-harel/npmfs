package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/g-harel/rejstry/internal/tarball"
	"github.com/g-harel/rejstry/internal/git"
	"github.com/g-harel/rejstry/internal/registry"
)

func v1Diff(w http.ResponseWriter, r *http.Request) {
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

	// Fetch package contents.
	pkg, err := registry.PackageContents(req.Registry, req.Package, req.Version)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR fetch package contents: %v", err)
		return
	}
	defer pkg.Close()

	// Fetch compare contents.
	cmp, err := registry.PackageContents(req.Registry, req.Package, req.Compare)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR fetch compare contents: %v", err)
		return
	}
	defer cmp.Close()

	// Create temporary working directory.
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR fetch compare contents: %v", err)
		return
	}

	// Initialize git repository.
	repo, err := git.Init(dir)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR git init %v", err)
		return
	}

	// Write package contents to directory.
	tarball.Extract(pkg, tarball.Downloader(dir))

	//
	err = repo.Add(".")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR git add %v: %v", req.Version, err)
		return
	}

	//
	err = repo.Commit(req.Version)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR git commit %v: %v", req.Version, err)
		return
	}

	// Write package contents to directory.
	tarball.Extract(cmp, tarball.Downloader(dir))

	//
	err = repo.Add(".")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR git add %v: %v", req.Compare, err)
		return
	}

	//
	err = repo.Commit(req.Compare)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR git commit %v: %v", req.Compare, err)
		return
	}

	//
	out, err := repo.DiffTree("HEAD", "HEAD~")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR git diff-tree: %v", err)
		return
	}

	//
	_, err = io.Copy(w, out)
	if err != nil {
		log.Printf("ERROR copy output: %v", err)
	}

	//
	err = os.RemoveAll(dir)
	if err != nil {
		log.Printf("ERROR cleanup: %v", err)
	}
}
