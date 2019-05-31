package diff

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

// Compare diffs the contents of two directories.
func Compare(a, b string) ([]*Patch, error) {
	// Create temporary working directory.
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return nil, fmt.Errorf("create temp dir: %v", err)
	}
	contentPath := path.Join(dir, "content")

	// Initialize git repository.
	_, err = execGit(dir, "init")
	if err != nil {
		return nil, err
	}

	// Move contents from "a" into repository.
	err = os.Rename(a, contentPath)
	if err != nil {
		return nil, fmt.Errorf("copy contents: %v", err)
	}

	// Commit version "a" to the repository.
	_, err = execGit(dir, "add", ".")
	if err != nil {
		return nil, err
	}
	_, err = execGit(dir, "commit", "-m", a)
	if err != nil {
		return nil, err
	}

	// Return contents from "a" to original path.
	err = os.Rename(contentPath, a)
	if err != nil {
		return nil, fmt.Errorf("copy contents: %v", err)
	}

	// Move contents from "b" into repository.
	err = os.Rename(b, contentPath)
	if err != nil {
		return nil, fmt.Errorf("copy contents: %v", err)
	}

	// Commit version "b" to the repository.
	_, err = execGit(dir, "add", ".")
	if err != nil {
		return nil, err
	}
	_, err = execGit(dir, "commit", "-m", b)
	if err != nil {
		return nil, err
	}

	// Compute diff between contents.
	out, err := execGit(dir, "diff-tree", "--patch", "-r", "--find-renames", "HEAD~", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("compute diff: %v", err)
	}

	// Return contents from "b" to original path.
	err = os.Rename(contentPath, b)
	if err != nil {
		return nil, fmt.Errorf("copy contents: %v", err)
	}

	// Clean up temporary directory.
	_ = os.RemoveAll(dir)

	// Parse output text.
	patches, err := parse(out)
	if err != nil {
		return nil, fmt.Errorf("parse output: %v", err)
	}

	return patches, nil
}
