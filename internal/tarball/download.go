package tarball

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

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
