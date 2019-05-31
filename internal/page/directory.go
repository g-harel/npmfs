package page

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/g-harel/rejstry/internal/registry"
	"github.com/g-harel/rejstry/internal/tarball"
)

func Directory(w http.ResponseWriter, r *http.Request, name, version, path string) {
	tmpl, err := template.ParseFiles(
		"templates/layout.html",
		"templates/pages/directory.html",
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

	parts, links := breakPath(path)
	context := &struct {
		Package     string
		Version     string
		Path        []string
		PathLinks   []string
		Directories []string
		Files       []string
	}{
		Package:     name,
		Version:     version,
		Path:        parts,
		PathLinks:   links,
		Directories: cleanup(dirs),
		Files:       cleanup(files),
	}

	err = tmpl.Execute(w, context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ERROR execute template: %v", err)
		return
	}
}
