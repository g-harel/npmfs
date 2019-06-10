package standard

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/g-harel/npmfs/internal/registry"
)

type stdReadHandler func(name string, contents io.Reader) error

// Helper to fetch the contents of a package and call the handler with each file.
func (c *Client) read(name, version string, handler stdReadHandler) error {
	client := &http.Client{Timeout: 7 * time.Second}

	// Fetch .tgz archive of package contents.
	url := fmt.Sprintf("https://%s/%s/-/%[2]s-%s.tgz", c.Host, name, version)
	response, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("request contents: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return registry.ErrNotFound
	}
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf(http.StatusText(response.StatusCode))
	}

	// Extract tarball from body.
	extractedBody, err := gzip.NewReader(response.Body)
	if err != nil {
		return fmt.Errorf("extract gzip data: %v", err)
	}

	// Extract package contents from tarball.
	tarball := tar.NewReader(extractedBody)
	for {
		header, err := tarball.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("advance to next file: %v", err)
		}

		if header.FileInfo().IsDir() {
			continue
		}

		err = handler(header.Name, tarball)
		if err != nil {
			return fmt.Errorf("handler error: %v", err)
		}
	}

	return nil
}
