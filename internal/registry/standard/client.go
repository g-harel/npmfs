package standard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/g-harel/npmfs/internal/registry"
	"github.com/g-harel/npmfs/internal/util"
)

// Client implements the registry.Client interface for standard registries.
type Client struct {
	Host string
}

var _ registry.Client = &Client{}

// Directory reads files and sub-directories at the given path.
func (c *Client) Directory(name, version, path string) ([]string, []string, error) {
	dirs := []string{}
	files := []string{}
	err := c.read(name, version, func(name string, contents io.Reader) error {
		filepath := strings.TrimPrefix(name, "package/")
		if strings.HasPrefix(filepath, path) {
			filepath := strings.TrimPrefix(filepath, path)
			pathParts := strings.Split(filepath, "/")
			if len(pathParts) == 1 {
				files = append(files, pathParts[0])
			} else {
				dirs = append(dirs, pathParts[0])
			}
		}
		return nil
	})
	if err == registry.ErrNotFound {
		return nil, nil, registry.ErrNotFound
	}
	if err != nil {
		return nil, nil, fmt.Errorf("read package contents: %v", err)
	}
	if len(dirs) == 0 && len(files) == 0 {
		return nil, nil, registry.ErrNotFound
	}

	return util.Unique(dirs), util.Unique(files), nil
}

// Download writes a package's contents to a temporary directory and returns its path.
func (c *Client) Download(name, version string) (string, error) {
	// Create temporary directory.
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", fmt.Errorf("create temp dir: %v", err)
	}

	// Write package contents to the target directory.
	err = c.read(name, version, func(name string, contents io.Reader) error {
		outPath := path.Join(dir, strings.TrimPrefix(name, "package"))

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
	})
	if err == registry.ErrNotFound {
		return "", registry.ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("read package contents: %v", err)
	}

	return dir, nil
}

// File reads a file's contents at the given path.
func (c *Client) File(name, version, path string) (string, error) {
	file := ""
	found := false
	err := c.read(name, version, func(name string, contents io.Reader) error {
		if !found && strings.TrimPrefix(name, "package/") == path {
			buf := new(bytes.Buffer)
			_, err := buf.ReadFrom(contents)
			if err != nil {
				return fmt.Errorf("ERROR copy contents: %v", err)
			}
			file = buf.String()
			found = true
		}
		return nil
	})
	if err == registry.ErrNotFound {
		return "", registry.ErrNotFound
	}
	if err != nil {
		return "", fmt.Errorf("read package contents: %v", err)
	}
	if !found {
		return "", registry.ErrNotFound
	}

	return file, nil
}

// Versions fetches all package versions from the registry.
func (c *Client) Versions(name string) ([]string, string, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	url := fmt.Sprintf("https://%s/%s", c.Host, name)
	response, err := client.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("request contents: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusNotFound {
		return nil, "", registry.ErrNotFound
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
