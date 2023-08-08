package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hiroara/drawin/job"
)

var downloadFailure = errors.New("Download failed.")

type Client struct {
	out Output
}

type Output interface {
	Add(j *job.Job, data []byte) error
	Check(j *job.Job) (bool, error)
	Prepare() error
}

func New(out Output) *Client {
	return &Client{out: out}
}

func Build(out Output) (*Client, error) {
	if err := out.Prepare(); err != nil {
		return nil, err
	}
	return New(out), nil
}

func (d *Client) Download(ctx context.Context, j *job.Job) error {
	ok, err := d.out.Check(j)
	if err != nil {
		return err
	}
	if ok {
		j.Action = job.CacheAction
		return nil
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

	if err := d.out.Add(j, body); err != nil {
		return err
	}
	j.Action = job.DownloadAction
	j.ContentLength = resp.ContentLength

	return nil
}
