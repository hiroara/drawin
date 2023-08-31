package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/hiroara/drawin"
	"github.com/hiroara/drawin/handler"
	"github.com/hiroara/drawin/job"
)

var downloadFailure = errors.New("Download failed.")

type Client struct {
	out      Output
	handlers []handler.Handler
	retry    *RetryConfig
}

type Output interface {
	Add(*drawin.Report, []byte) error
	Get(*job.Job) (*drawin.Report, error)
	Initialize() error
}

func New(out Output, handlers []handler.Handler, opts ...Option) *Client {
	if handlers == nil {
		handlers = handler.DefaultHandlers
	}
	cli := &Client{out: out, handlers: handlers}
	for _, opt := range opts {
		opt(cli)
	}
	return cli
}

func Build(out Output, handlers []handler.Handler, opts ...Option) (*Client, error) {
	if err := out.Initialize(); err != nil {
		return nil, err
	}
	return New(out, handlers, opts...), nil
}

func (cli *Client) Download(ctx context.Context, j *job.Job) (*drawin.Report, error) {
	rep, err := cli.out.Get(j)
	if err != nil {
		return nil, err
	}

	h, err := cli.selectHandler(j)
	if err != nil {
		return nil, err
	}

	if cli.useCache(rep) {
		rep.Result = drawin.CachedResult
		return rep, nil
	}

	bs, err := h.Get(ctx, j)
	if err != nil {
		rep := drawin.FailedReport(j, err, !h.ShouldRetry(err))

		if err := cli.out.Add(rep, bs); err != nil {
			return nil, err
		}

		return rep, nil
	}

	rep = drawin.DownloadedReport(j, int64(len(bs)))

	if err := cli.out.Add(rep, bs); err != nil {
		return nil, err
	}

	return rep, nil
}

var errNoMatchingHandler = errors.New("no matching handler is found")

func (cli *Client) selectHandler(j *job.Job) (handler.Handler, error) {
	for _, h := range cli.handlers {
		if h.Match(j) {
			return h, nil
		}
	}
	return nil, fmt.Errorf("%w for job: %s (URL: %s)", errNoMatchingHandler, j.Name, j.URL)
}

func (cli *Client) useCache(rep *drawin.Report) bool {
	if rep == nil {
		return false
	}
	return !cli.retry.shouldRetry(rep)
}
