package registry

import "io"

// Client defines the required interface for a registry.
type Client interface {
	// Archive writes a zip archive of the package contents at path to out.
	Archive(name, version, path string, out io.Writer) (err error)
	// Directory reads files and sub-directories at the given path.
	Directory(name, version, path string) (dirs, files []string, err error)
	// Download writes a package's contents to a temporary directory and returns its path.
	Download(name, version string) (dir string, err error)
	// File reads a file's contents at the given path.
	File(name, version, path string) (file string, err error)
	// Versions fetches all package versions from the registry.
	Versions(name string) (versions []string, latest string, err error)
}
