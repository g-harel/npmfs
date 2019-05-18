package internal

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
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

// Downloader generates a ExtractHandler which writes contents to a directory.
func Downloader(dir string) ExtractHandler {
	return func(name string, contents io.Reader) error {
		outPath := filepath.Join(dir, name)

		err := os.MkdirAll(path.Dir(outPath), os.ModePerm)
		if err != nil {
			return fmt.Errorf("create file output path: %v", err)
		}

		outFile, err := os.Create(outPath)
		if err != nil {
			return fmt.Errorf("create output file: %v", err)
		}

		_, err = io.Copy(outFile, contents)
		if err != nil {
			return fmt.Errorf("write file contents: %v", err)
		}

		err = outFile.Close()
		if err != nil {
			return fmt.Errorf("close file: %v", err)
		}

		return nil
	}
}
