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
	ShouldRetry(error) bool
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

func (cli *Client) Download(ctx context.Context, j *job.Job) (*reporter.Report, error) {
	rep, err := cli.out.Get(j)
	if err != nil {
		return nil, err
	}

	h, err := cli.selectHandler(j)
	if err != nil {
		return nil, err
	}

	if cli.useCache(rep) {
		rep.Result = reporter.Cached
		return rep, nil
	}

	bs, err := h.Get(ctx, j)
	if err != nil {
		rep := reporter.FailedReport(j, err, !h.ShouldRetry(err))

		if err := cli.out.Add(rep, bs); err != nil {
			return nil, err
		}

		return rep, nil
	}

	rep = reporter.DownloadedReport(j, int64(len(bs)))

	if err := cli.out.Add(rep, bs); err != nil {
		return nil, err
	}

	return rep, nil
}

var errNoMatchingHandler = errors.New("no matching handler is found")

func (cli *Client) selectHandler(j *job.Job) (Handler, error) {
	for _, h := range cli.handlers {
		if h.Match(j) {
			return h, nil
		}
	}
	return nil, fmt.Errorf("%w for job: %s (URL: %s)", errNoMatchingHandler, j.Name, j.URL)
}

func (cli *Client) useCache(rep *reporter.Report) bool {
	if rep == nil {
		return false
	}
	if rep.Result != reporter.Failed {
		return true
	}
	return rep.Failure.Permanent
}
