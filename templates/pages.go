package templates

import (
	"github.com/g-harel/npmfs/internal/diff"
)

func PageCompare(name, versionA, versionB string, patches []*diff.Patch) *Renderer {
	return &Renderer{
		filenames: []string{
			"templates/layout.html",
			"templates/logo.html",
			"templates/pages/compare.html",
		},
		context: struct {
			Package  string
			VersionA string
			VersionB string
			Patches  []*diff.Patch
		}{
			Package:  name,
			VersionA: versionA,
			VersionB: versionB,
			Patches:  patches,
		},
	}
}

func PageDirectory(name, version string, path, links, dirs, files []string) *Renderer {
	return &Renderer{
		filenames: []string{
			"templates/layout.html",
			"templates/logo.html",
			"templates/pages/directory.html",
		},
		context: struct {
			Package     string
			Version     string
			Path        []string
			PathLinks   []string
			Directories []string
			Files       []string
		}{
			Package:     name,
			Version:     version,
			Path:        path,
			PathLinks:   links,
			Directories: dirs,
			Files:       files,
		},
	}
}

func PageFile(name, version string, path, links, lines []string) *Renderer {
	return &Renderer{
		filenames: []string{
			"templates/layout.html",
			"templates/logo.html",
			"templates/pages/file.html",
		},
		context: struct {
			Package   string
			Version   string
			Path      []string
			PathLinks []string
			Lines     []string
		}{
			Package:   name,
			Version:   version,
			Path:      path,
			PathLinks: links,
			Lines:     append([]string{""}, lines...),
		},
	}
}

func PageVersions(name, latest, disabled string, versions []string) *Renderer {
	return &Renderer{
		filenames: []string{
			"templates/layout.html",
			"templates/logo.html",
			"templates/pages/versions.html",
		},
		context: struct {
			Package  string
			Latest   string
			Disabled string
			Versions []string
		}{
			Package:  name,
			Latest:   latest,
			Disabled: disabled,
			Versions: versions,
		},
	}
}
