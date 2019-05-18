package rejstry

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

type Package struct {
	Tags struct {
		Latest string `json:"latest"`
	} `json:"dist-tags"`
	Versions map[string]struct {
		Dist struct {
			Tarball string `json:"tarball"`
		} `json:"dist"`
	} `json:"versions"`
}

// Download writes the package contents to a temporary directory and returns its path.
// The tarball contents are expected to have been compressed using gzip.
func (p *Package) Download(version string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	if _, ok := p.Versions[version]; !ok {
		return "", fmt.Errorf("unrecognized version")
	}
	println(p.Versions[version].Dist.Tarball)
	response, err := client.Get(p.Versions[version].Dist.Tarball)
	if err != nil {
		return "", fmt.Errorf("fetch package contents: %v", err)
	}
	defer response.Body.Close()

	outputDir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", fmt.Errorf("create output directory: %v", err)
	}

	contents, err := gzip.NewReader(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tarball := tar.NewReader(contents)
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
