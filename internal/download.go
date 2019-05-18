package internal

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

// DownloadPackage writes the package contents to a temporary directory and returns its path.
// The registry's tarball contents are expected to be compressed using gzip.
func DownloadPackage(registry, name, version string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	url := fmt.Sprintf("https://%s/%s/-/%[2]s-%s.tgz", registry, name, version)
	response, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("request contents: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf(http.StatusText(response.StatusCode))
	}

	outputDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", fmt.Errorf("create output directory: %v", err)
	}

	extractedBody, err := gzip.NewReader(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tarball := tar.NewReader(extractedBody)
	for {
		header, err := tarball.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("read next file: %v", err)
		}

		outPath := filepath.Join(outputDir, header.Name)

		err = os.MkdirAll(path.Dir(outPath), os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("create file output path: %v", err)
		}

		outFile, err := os.Create(outPath)
		if err != nil {
			return "", fmt.Errorf("create output file: %v", err)
		}

		_, err = io.Copy(outFile, tarball)
		if err != nil {
			return "", fmt.Errorf("write file contents: %v", err)
		}

		outFile.Close()
	}

	return outputDir, nil
}
