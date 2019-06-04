package registry

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type standardRegistry struct {
	host string
}

// PackageContents fetches the data for a package's contents.
func (s *standardRegistry) PackageContents(name, version string) (io.ReadCloser, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	url := fmt.Sprintf("https://%s/%s/-/%[2]s-%s.tgz", s.host, name, version)
	response, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request contents: %v", err)
	}

	if response.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(http.StatusText(response.StatusCode))
	}

	return response.Body, nil
}

// PackageVersions fetches all package versions from the registry.
// Returned list of versions is sorted in descending order.
// Latest version is also returned in the second position.
func (s *standardRegistry) PackageVersions(name string) ([]string, string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	url := fmt.Sprintf("https://%s/%s", s.host, name)
	response, err := client.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("request contents: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, "", ErrNotFound
	}
	if response.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf(http.StatusText(response.StatusCode))
	}

	data := &struct {
		Versions map[string]interface{} `json:"versions"`
		Tags     struct {
			Latest string `json:"latest"`
		} `json:"dist-tags"`
	}{}

	err = json.NewDecoder(response.Body).Decode(data)
	if err != nil {
		return nil, "", fmt.Errorf("decode response body: %v", err)
	}

	versions := make([]string, len(data.Versions))
	count := 0
	for version := range data.Versions {
		versions[count] = version
		count++
	}

	return versions, data.Tags.Latest, nil
}
