package templates

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"github.com/g-harel/rejstry/internal/diff"
	"github.com/g-harel/rejstry/internal/registry"
	"github.com/g-harel/rejstry/internal/tarball"
)

func Compare(w http.ResponseWriter, r *http.Request, name, versionA, versionB string) {
	tmpl, err := template.ParseFiles(
		"templates/layout.html",
		"templates/pages/compare.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ERROR parse template: %v", err)
		return
	}

	type downloadedDir struct {
		version string
		dir     string
		err     error
	}

	dirChan := make(chan downloadedDir)
	for _, version := range []string{versionA, versionB} {
		go func(v string) {
			// Create temporary working directory.
			dir, err := ioutil.TempDir("", "")
			if err != nil {
				dirChan <- downloadedDir{v, "", fmt.Errorf("create temp dir: %v", err)}
				return
			}

			// Fetch package contents for given version.
			pkg, err := registry.PackageContents("registry.npmjs.com", name, v)
			if err != nil {
				dirChan <- downloadedDir{v, "", fmt.Errorf("fetch package: %v", err)}
				return
			}
			defer pkg.Close()

			// Write package contents to directory.
			tarball.Extract(pkg, tarball.Downloader(dir))

			dirChan <- downloadedDir{v, dir, nil}
		}(version)
	}

	dirs := map[string]string{}
	for i := 0; i < 2; i++ {
		dir := <-dirChan
		if dir.err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Printf("ERROR download package '%v': %v", dir.version, err)
			return
		}
		dirs[dir.version] = dir.dir
	}

	patches, err := diff.Compare(path.Join(dirs[versionA], "package"), path.Join(dirs[versionB], "package"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Printf("ERROR compare package contents: %v", err)
		return
	}

	context := &struct {
		Package  string
		VersionA string
		VersionB string
		Patches  []*diff.Patch
	}{
		Package:  name,
		VersionA: versionA,
		VersionB: versionB,
		Patches:  patches,
	}

	err = tmpl.Execute(w, context)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ERROR execute template: %v", err)
		return
	}
}
