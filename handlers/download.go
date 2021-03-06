package handlers

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/g-harel/npmfs/internal/registry"
	"github.com/g-harel/npmfs/templates"
	"github.com/gorilla/mux"
	"golang.org/x/xerrors"
)

// DownloadFile handler serves a file from the package contents.
func DownloadFile(client registry.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		version := vars["version"]
		path := vars["path"]

		filename := filepath.Base(path)

		w.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(filename)))
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%v", filename))

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
		io.WriteString(w, file)
	}
}

// DownloadDir handler serves a zip archive of the package contents.
func DownloadDir(client registry.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		version := vars["version"]
		path := vars["path"]

		filename := fmt.Sprintf("%v-%v-%v", name, version, strings.ReplaceAll(path, "/", "-"))
		filename = strings.TrimSuffix(filename, "-")
		filename += ".zip"

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%v", filename))

		err := client.Archive(name, version, path, w)
		if err != nil {
			var registryErr *registry.Error
			if xerrors.As(err, &registryErr) {
				templates.PageError(registryErr.StatusCode, registryErr.Error()).Handler(w, r)
				return
			}
			templates.PageError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)).Handler(w, r)
			log.Printf("ERROR fetch package archive: %v", err)
			return
		}
	}
}
