package tarball

import (
	"fmt"
	"io"
	"os"
	"path"
)

// Downloader generates a ExtractHandler which writes contents to a directory.
// It accepts a function to pick the output path from the name.
func Downloader(picker func(string) string) ExtractHandler {
	return func(name string, contents io.Reader) error {
		outPath := picker(name)

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
