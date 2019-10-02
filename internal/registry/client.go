package registry

import "io"

// Client defines the required interface for a registry.
type Client interface {
	Archive(name, version string, out io.Writer) (err error)
	Directory(name, version, path string) (dirs, files []string, err error)
	Download(name, version string) (dir string, err error)
	File(name, version, path string) (file string, err error)
	Versions(name string) (versions []string, latest string, err error)
}
