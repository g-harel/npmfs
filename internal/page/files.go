package page

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/g-harel/rejstry/internal/registry"
	"github.com/g-harel/rejstry/internal/tarball"
)

func Files(w http.ResponseWriter, r *http.Request, name, version string) {
	tmpl, err := template.ParseFiles(
		"templates/layout.html",
		"templates/pages/files.html",
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

	// Extract files from package contents.
	fileList := []string{}
	err = tarball.Extract(pkg, func(name string, contents io.Reader) error {
		fileList = append(fileList, strings.TrimPrefix(name, "package/"))
		return nil
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR extract files from package contents: %v", err)
		return
	}
	sort.Strings(fileList)

	type pkgFile struct {
		Depth    int
		Path     string
		Children map[string]pkgFile
	}

	// Create tree structure from full file paths.
	files := map[string]pkgFile{}
	for _, filePath := range fileList {
		temp := files
		for i, part := range strings.Split(filePath, "/") {
			if part == "" {
				break
			}
			if _, ok := temp[part]; !ok {
				temp[part] = pkgFile{
					Depth:    i,
					Path:     filePath,
					Children: map[string]pkgFile{},
				}
			}
			temp = temp[part].Children
		}
	}

	context := &struct {
		Package string
		Version string
		Files   map[string]pkgFile
	}{
		Package: name,
		Version: version,
		Files:   files,
	}

	err = tmpl.Execute(w, context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ERROR execute template: %v", err)
		return
	}
}
