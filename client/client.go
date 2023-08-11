package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/reporter"
)

var downloadFailure = errors.New("Download failed.")

type Client struct {
	out      Output
	handlers []Handler
}

type Handler interface {
	Match(*job.Job) bool
	Get(context.Context, *job.Job) ([]byte, error)
}

var DefaultHandlers = []Handler{&HTTPHandler{client: http.DefaultClient}}

type Output interface {
	Add(*reporter.Report, []byte) error
	Get(*job.Job) (*reporter.Report, error)
	Initialize() error
}

func New(out Output, opts ...Option) *Client {
	cli := &Client{out: out, handlers: DefaultHandlers}
	for _, opt := range opts {
		opt(cli)
	}
	return cli
}

func Build(out Output, opts ...Option) (*Client, error) {
	if err := out.Initialize(); err != nil {
		return nil, err
	}
	return New(out, opts...), nil
}

var ErrUnexpectedResponseStatus = errors.New("received unexpected HTTP response status code")

func (d *Client) Download(ctx context.Context, j *job.Job) (*reporter.Report, error) {
	rep, err := d.out.Get(j)
	if err != nil {
		return nil, err
	}
	if rep != nil {
		rep.Result = reporter.Cached
		return rep, nil
	}

	bs, err := d.download(ctx, j)
	if err != nil {
		rep := reporter.FailedReport(j, err)

		if err := d.out.Add(rep, bs); err != nil {
			return nil, err
		}

		return rep, nil
	}

	rep = reporter.DownloadedReport(j, int64(len(bs)))

	if err := d.out.Add(rep, bs); err != nil {
		return nil, err
	}

	return rep, nil
}

var errNoMatchingHandler = errors.New("no matching handler is found")

func (d *Client) download(ctx context.Context, j *job.Job) ([]byte, error) {
	h := d.findHandler(j)
	if h == nil {
		return nil, fmt.Errorf("%w for job: %s (URL: %s)", errNoMatchingHandler, j.Name, j.URL)
	}
	return h.Get(ctx, j)
}

func (d *Client) findHandler(j *job.Job) Handler {
	for _, h := range d.handlers {
		if h.Match(j) {
			return h
		}
	}
	return nil
}
