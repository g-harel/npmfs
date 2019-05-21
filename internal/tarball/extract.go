package tarball

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
)

// ExtractHandler is invoked for each file with it's name and contents.
type ExtractHandler func(name string, contents io.Reader) error

// Extract extracts the contents of a gzipped tarball using the provided handler.
func Extract(source io.Reader, handler ExtractHandler) error {
	extractedBody, err := gzip.NewReader(source)
	if err != nil {
		return fmt.Errorf("extract gzip data: %v", err)
	}

	tarball := tar.NewReader(extractedBody)
	for {
		header, err := tarball.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("advance to next file: %v", err)
		}

		err = handler(header.Name, tarball)
		if err != nil {
			return fmt.Errorf("handler error: %v", err)
		}
	}

	return nil
}
