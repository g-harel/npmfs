package internal

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// FetchPackage fetches the data for a package's contents.
func FetchPackage(registry, name, version string) (io.ReadCloser, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	url := fmt.Sprintf("https://%s/%s/-/%[2]s-%s.tgz", registry, name, version)
	response, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request contents: %v", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(http.StatusText(response.StatusCode))
	}

	return response.Body, nil
}
