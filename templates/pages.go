package templates

import (
	"github.com/g-harel/npmfs/internal/diff"
)

// PageHome returns a renderer for the home page.
func PageHome() *Renderer {
	return &Renderer{
		filenames: []string{
			"templates/layout.html",
			"templates/logo.html",
			"templates/pages/home.html",
		},
		context: nil,
	}
}

// PageCompare returns a renderer for the compare page.
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

// PageDirectory returns a renderer for the directory page.
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

// PageFile returns a renderer for the file page.
func PageFile(name, version, size string, path, links, lines []string) *Renderer {
	return &Renderer{
		filenames: []string{
			"templates/layout.html",
			"templates/logo.html",
			"templates/pages/file.html",
		},
		context: struct {
			Package   string
			Version   string
			Size      string
			Path      []string
			PathLinks []string
			Lines     []string
		}{
			Package:   name,
			Version:   version,
			Size:      size,
			Path:      path,
			PathLinks: links,
			Lines:     append([]string{""}, lines...),
		},
	}
}

// PageVersions returns a renderer for the versions page.
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

// PageError returns a renderer for a generic error page.
func PageError(status int, info string) *Renderer {
	return &Renderer{
		statusCode: status,
		filenames: []string{
			"templates/layout.html",
			"templates/logo.html",
			"templates/pages/error.html",
		},
		context: struct {
			Status int
			Info   string
		}{
			Status: status,
			Info:   info,
		},
	}
}
