package client

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/reporter"
)

type DirectoryOutput struct {
	dir string
}

func NewDirectory(dir string) *DirectoryOutput {
	return &DirectoryOutput{dir: dir}
}

func (out *DirectoryOutput) Add(rep *reporter.Report, data []byte) error {
	if data == nil {
		// Does not store anything when data is missing
		return nil
	}
	return os.WriteFile(out.fullpath(rep.Name), data, 0644)
}

func (out *DirectoryOutput) Get(j *job.Job) (*reporter.Report, error) {
	stat, err := os.Stat(out.fullpath(j.Name))
	if err == nil { // File exists
		return reporter.CachedReport(j, stat.Size()), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	return nil, err
}

func (out *DirectoryOutput) Initialize() error {
	return os.MkdirAll(out.dir, 0755)
}

func (out *DirectoryOutput) fullpath(name string) string {
	return filepath.Join(out.dir, name)
}
