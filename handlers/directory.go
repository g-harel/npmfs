package handlers

import (
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/g-harel/rejstry/internal/paths"
	"github.com/g-harel/rejstry/internal/registry"
	"github.com/g-harel/rejstry/internal/tarball"
	"github.com/g-harel/rejstry/templates"
	"github.com/gorilla/mux"
)

func Directory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	version := vars["version"]
	path := vars["path"]

	// Fetch package contents.
	pkg, err := registry.PackageContents("registry.npmjs.com", name, version)
	if err == registry.ErrNotFound {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR fetch package contents: %v", err)
		return
	}
	defer pkg.Close()

	// Extract files and directories at the given path.
	dirs := []string{}
	files := []string{}
	err = tarball.Extract(pkg, func(name string, contents io.Reader) error {
		filePath := strings.TrimPrefix(name, "package/")
		if strings.HasPrefix(filePath, path) {
			filePath := strings.TrimPrefix(filePath, path)
			pathParts := strings.Split(filePath, "/")
			if len(pathParts) == 1 {
				files = append(files, pathParts[0])
			} else {
				dirs = append(dirs, pathParts[0])
			}
		}
		return nil
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR extract files from package contents: %v", err)
		return
	}
	if len(dirs) == 0 && len(files) == 0 {
		http.NotFound(w, r)
		return
	}

	// Sort and de-duplicate input slice.
	cleanup := func(s []string) []string {
		m := map[string]interface{}{}
		for _, item := range s {
			m[item] = true
		}
		out := []string{}
		for key := range m {
			out = append(out, key)
		}
		sort.Strings(out)
		return out
	}

	dirs = cleanup(dirs)
	files = cleanup(files)
	parts, links := paths.BreakRelative(path)

	templates.PageDirectory(name, version, parts, links, dirs, files).Render(w)
}
