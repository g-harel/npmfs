package registry

import (
	"errors"
	"io"
	"net/http"
)

var ErrNotFound = errors.New(http.StatusText(http.StatusNotFound))

type Registry interface {
	PackageVersions(name string) ([]string, string, error)
	PackageContents(name, version string) (io.ReadCloser, error)
}

var NPM Registry = &standardRegistry{"registry.npmjs.com"}
var Yarn Registry = &standardRegistry{"registry.yarnpkg.com"}
var Open Registry = &standardRegistry{"npm.open-registry.dev"}
