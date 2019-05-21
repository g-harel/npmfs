package git

import (
	"fmt"
	"io"
)

type Repository struct {
	dir string
}

func Init(dir string) (*Repository, error) {
	_, err := execGit(dir, "init")
	if err != nil {
		return nil, fmt.Errorf("git init: %v", err)
	}

	return &Repository{dir}, nil
}

func (r *Repository) Add(path string) error {
	_, err := execGit(r.dir, "add", path)
	return err
}

func (r *Repository) Commit(message string) error {
	_, err := execGit(r.dir, "commit", "-m", message)
	return err
}

func (r *Repository) DiffTree(a, b string) (io.Reader, error) {
	// return execGit(r.dir, "diff-tree", "--name-status", "-r", "--find-renames", a, b)
	return execGit(r.dir, "diff-tree", "--patch", "-r", "--find-renames", a, b)
}
