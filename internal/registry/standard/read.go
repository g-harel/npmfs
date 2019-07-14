package standard

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/g-harel/npmfs/internal/registry"
	"golang.org/x/xerrors"
)

type stdReadHandler func(name string, contents io.Reader) error

// Helper to fetch the contents of a package and call the handler with each file.
func (c *Client) read(name, version string, handler stdReadHandler) error {
	client := &http.Client{Timeout: 4 * time.Second}

	// Fetch .tgz archive of package contents.
	url := fmt.Sprintf("https://%s/%s/-/%[2]s-%s.tgz", c.Host, name, version)
	response, err := client.Get(url)
	if os.IsTimeout(err) {
		log.Printf("ERROR standard registry: read: timeout (%v)", url)
		return registry.ErrGatewayTimeout
	}
	if err != nil {
		return xerrors.Errorf("request contents: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		log.Printf("ERROR standard registry: read: not found (%v)", url)
		return registry.ErrNotFound
	}
	if response.StatusCode != http.StatusOK {
		log.Printf("ERROR standard registry: read: unexpected status code (%v): %v", url, response.StatusCode)
		return registry.ErrBadGateway
	}

	// Extract tarball from body.
	extractedBody, err := gzip.NewReader(response.Body)
	if err != nil {
		return xerrors.Errorf("extract gzip data: %w", err)
	}

	// Extract package contents from tarball.
	tarball := tar.NewReader(extractedBody)
	for {
		header, err := tarball.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return xerrors.Errorf("advance to next file: %w", err)
		}

		if header.FileInfo().IsDir() {
			continue
		}

		err = handler(header.Name, tarball)
		if err != nil {
			return xerrors.Errorf("handler error: %w", err)
		}
	}

	return nil
}
