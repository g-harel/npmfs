package mock

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/g-harel/npmfs/internal/registry"
)

// Registry is a mock implementation of the registry.Registry interface.
type Registry struct {
	Contents           map[string]map[string]string
	Latest             string
	PackageVersionsErr error
	PackageContentsErr error
}

var _ registry.Registry = &Registry{}

// PackageVersions returns all versions listed in the contents and the specified latest value.
// Package name is ignored.
func (ry *Registry) PackageVersions(name string) ([]string, string, error) {
	if ry.PackageVersionsErr != nil {
		return nil, "", ry.PackageVersionsErr
	}

	versions := []string{}
	for v := range ry.Contents {
		versions = append(versions, v)
	}
	return versions, ry.Latest, ry.PackageVersionsErr
}

// PackageContents returns a gzipped tarball of the contents of the specified version.
// Package name is ignored.
func (ry *Registry) PackageContents(name, version string) (io.ReadCloser, error) {
	if ry.PackageContentsErr != nil {
		return nil, ry.PackageContentsErr
	}

	versionContents, ok := ry.Contents[version]
	if !ok {
		return nil, registry.ErrNotFound
	}

	// Create gzip/tar writer chain.
	buf := &bytes.Buffer{}
	gw := gzip.NewWriter(buf)
	tw := tar.NewWriter(gw)

	// Write version's contents.
	for name, content := range versionContents {
		hdr := &tar.Header{
			Name: name,
			Mode: 0600,
			Size: int64(len(content)),
		}
		err := tw.WriteHeader(hdr)
		if err != nil {
			return nil, fmt.Errorf("write file header: %v", err)
		}
		_, err = tw.Write([]byte(content))
		if err != nil {
			return nil, fmt.Errorf("write file content: %v", err)
		}
	}
	err := tw.Close()
	if err != nil {
		return nil, fmt.Errorf("close tarball writer: %v", err)
	}
	err = gw.Close()
	if err != nil {
		return nil, fmt.Errorf("close gzip writer: %v", err)
	}

	return ioutil.NopCloser(buf), ry.PackageContentsErr
}
