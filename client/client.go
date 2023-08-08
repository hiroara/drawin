package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/reporter"
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

var ErrUnexpectedResponseStatus = errors.New("received unexpected HTTP response status code")

func (d *Client) Download(ctx context.Context, j *job.Job) (*reporter.Report, error) {
	ok, err := d.out.Check(j)
	if err != nil {
		return nil, err
	}
	if ok {
		return reporter.CachedReport(j), nil
	}

	resp, err := http.Get(j.URL)
	if err != nil {
		return reporter.FailedReport(j, err), nil
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return reporter.FailedReport(j, fmt.Errorf("%w: %d", ErrUnexpectedResponseStatus, resp.StatusCode)), nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return reporter.FailedReport(j, err), nil
	}

	if err := d.out.Add(j, body); err != nil {
		return nil, err
	}

	return reporter.DownloadedReport(j, resp.ContentLength), nil
}
