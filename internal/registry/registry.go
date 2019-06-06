package registry

import (
	"errors"
	"io"
	"net/http"
)

// ErrNotFound is returned when a package and version combination is not found.
var ErrNotFound = errors.New(http.StatusText(http.StatusNotFound))

// Registry defines the required interface for an external registry.
type Registry interface {
	PackageVersions(name string) ([]string, string, error)
	PackageContents(name, version string) (io.ReadCloser, error)
}

// Exported known public registries.
var (
	NPM  Registry = &standardRegistry{"registry.npmjs.com"}
	Yarn Registry = &standardRegistry{"registry.yarnpkg.com"}
	Open Registry = &standardRegistry{"npm.open-registry.dev"}
)
