package mock

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/g-harel/npmfs/internal/registry"
	"github.com/g-harel/npmfs/internal/util"
)

// Client is a mock implementation of the registry.Client interface.
type Client struct {
	Contents     map[string]map[string]string
	Latest       string
	DirectoryErr error
	DownloadErr  error
	FileErr      error
	VersionsErr  error
}

var _ registry.Client = &Client{}

// Directory lists all the sub-directories and files at the given path in the mocked contents.
// Package name is ignored.
func (c *Client) Directory(name, version, path string) ([]string, []string, error) {
	versionContents, ok := c.Contents[version]
	if !ok {
		return nil, nil, registry.ErrNotFound
	}

	dirs := []string{}
	files := []string{}
	for filepath := range versionContents {
		if strings.HasPrefix(filepath, path) {
			filepath := strings.TrimPrefix(filepath, path)
			pathParts := strings.Split(filepath, "/")
			if len(pathParts) == 1 {
				files = append(files, pathParts[0])
			} else {
				dirs = append(dirs, pathParts[0])
			}
		}
	}
	if len(dirs) == 0 && len(files) == 0 {
		return nil, nil, registry.ErrNotFound
	}

	return util.Unique(dirs), util.Unique(files), nil
}

// Download writes the mocked contents to a temporary directory.
// Package name is ignored.
func (c *Client) Download(name, version string) (string, error) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", fmt.Errorf("create temp dir: %v", err)
	}

	versionContents, ok := c.Contents[version]
	if !ok {
		return "", registry.ErrNotFound
	}

	for filepath, contents := range versionContents {
		outPath := path.Join(dir, filepath, "package")

		err := os.MkdirAll(path.Dir(outPath), os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("create file output path: %v", err)
		}

		outFile, err := os.Create(outPath)
		if err != nil {
			return "", fmt.Errorf("create output file: %v", err)
		}

		_, err = io.Copy(outFile, strings.NewReader(contents))
		if err != nil {
			return "", fmt.Errorf("write file contents: %v", err)
		}

		err = outFile.Close()
		if err != nil {
			return "", fmt.Errorf("close file: %v", err)
		}
	}

	return dir, c.DownloadErr
}

// File reads a file's contents at the given path in the mocked contents.
// Package name is ignored.
func (c *Client) File(name, version, path string) (string, error) {
	versionContents, ok := c.Contents[version]
	if !ok {
		return "", registry.ErrNotFound
	}

	fileContents, ok := versionContents[path]
	if !ok {
		return "", registry.ErrNotFound
	}

	return fileContents, c.FileErr
}

// Versions returns all versions listed in the contents and the specified latest value.
// Package name is ignored.
func (c *Client) Versions(name string) ([]string, string, error) {
	versions := []string{}
	for v := range c.Contents {
		versions = append(versions, v)
	}

	return versions, c.Latest, c.VersionsErr
}
