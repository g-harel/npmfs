package handlers

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/g-harel/npmfs/internal/paths"
	"github.com/g-harel/npmfs/internal/registry"
	"github.com/g-harel/npmfs/internal/tarball"
	"github.com/g-harel/npmfs/templates"
	"github.com/gorilla/mux"
)

// File handler displays a file view of package contents at the provided path.
func File(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	version := vars["version"]
	path := vars["path"]

	// Fetch package contents.
	pkg, err := registry.NPM.PackageContents(name, version)
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

	// Find file contents to use in response.
	// Contents must be written to a buffer to be used in a template.
	var file *bytes.Buffer
	err = tarball.Extract(pkg, func(name string, contents io.Reader) error {
		if strings.TrimPrefix(name, "package/") == path {
			file = new(bytes.Buffer)
			_, err := file.ReadFrom(contents)
			if err != nil {
				log.Printf("ERROR copy contents: %v", err)
			}
		}
		return nil
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR extract files from package contents: %v", err)
		return
	}
	if file == nil {
		http.NotFound(w, r)
		return
	}

	parts, links := paths.BreakRelative(path)
	lines := strings.Split(file.String(), "\n")

	// Render page template.
	templates.PageFile(name, version, parts, links, lines).Render(w)
}
