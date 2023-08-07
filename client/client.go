package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/hiroara/drawin/job"
)

var downloadFailure = errors.New("Download failed.")

type Client struct {
	dir string
}

func New(dir string) *Client {
	return &Client{dir: dir}
}

func (d *Client) CreateDir() error {
	return os.MkdirAll(d.dir, 0755)
}

func (d *Client) Download(ctx context.Context, j *job.Job) error {
	p := d.fullpath(j.Name)

	_, err := os.Stat(p)
	if err == nil { // File exists
		j.Action = job.CacheAction
		return nil // Bypass
	}

	resp, err := http.Get(j.URL)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("%w Unexpected response status code: %d", downloadFailure, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := store(p, body); err != nil {
		return err
	}
	j.Action = job.DownloadAction
	j.ContentLength = resp.ContentLength

	return nil
}

func store(p string, data []byte) error {
	return os.WriteFile(p, data, 0644)
}

func (d *Client) fullpath(name string) string {
	return path.Join(d.dir, name)
}
