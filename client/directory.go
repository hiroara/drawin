package client

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/hiroara/drawin/job"
)

type DirectoryOutput struct {
	dir string
}

func NewDirectory(dir string) *DirectoryOutput {
	return &DirectoryOutput{dir: dir}
}

func (out *DirectoryOutput) Add(j *job.Job, data []byte) error {
	return os.WriteFile(out.fullpath(j.Name), data, 0644)
}

func (out *DirectoryOutput) Check(j *job.Job) (bool, error) {
	_, err := os.Stat(out.fullpath(j.Name))
	if err == nil { // File exists
		return true, nil // Bypass
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func (out *DirectoryOutput) Prepare() error {
	return os.MkdirAll(out.dir, 0755)
}

func (out *DirectoryOutput) fullpath(name string) string {
	return filepath.Join(out.dir, name)
}
