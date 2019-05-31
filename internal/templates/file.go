package templates

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/g-harel/rejstry/internal/registry"
	"github.com/g-harel/rejstry/internal/tarball"
)

func File(w http.ResponseWriter, r *http.Request, name, version, path string) {
	tmpl, err := template.ParseFiles(
		"templates/layout.html",
		"templates/pages/file.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ERROR parse template: %v", err)
		return
	}

	// Fetch package contents.
	pkg, err := registry.PackageContents("registry.npmjs.com", name, version)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR fetch package contents: %v", err)
		return
	}
	defer pkg.Close()

	// Find file contents to use in response.
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

	parts, links := breakPath(path)
	context := &struct {
		Package   string
		Version   string
		Path      []string
		PathLinks []string
		Lines     []string
	}{
		Package:   name,
		Version:   version,
		Path:      parts,
		PathLinks: links,
		Lines:     strings.Split("\n"+file.String(), "\n"),
	}

	err = tmpl.Execute(w, context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ERROR execute template: %v", err)
		return
	}
}
