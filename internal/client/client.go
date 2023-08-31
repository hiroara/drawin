package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/hiroara/drawin/handler"
	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/output"
	"github.com/hiroara/drawin/report"
	"github.com/hiroara/drawin/retry"
)

var downloadFailure = errors.New("Download failed.")

type Client struct {
	out      output.Output
	handlers []handler.Handler
	retry    *retry.RetryConfig
}

func New(out output.Output, handlers []handler.Handler, retry *retry.RetryConfig) *Client {
	if handlers == nil {
		handlers = handler.DefaultHandlers
	}
	return &Client{out: out, handlers: handlers, retry: retry}
}

func Build(out output.Output, handlers []handler.Handler, retry *retry.RetryConfig) (*Client, error) {
	if err := out.Initialize(); err != nil {
		return nil, err
	}
	return New(out, handlers, retry), nil
}

func (cli *Client) Download(ctx context.Context, j *job.Job) (*report.Report, error) {
	rep, err := cli.out.Get(j)
	if err != nil {
		return nil, err
	}

	h, err := cli.selectHandler(j)
	if err != nil {
		return nil, err
	}

	if cli.useCache(rep) {
		rep.Result = report.CachedResult
		return rep, nil
	}

	bs, err := h.Get(ctx, j)
	if err != nil {
		rep := report.Failed(j, err, !h.ShouldRetry(err))

		if err := cli.out.Add(rep, bs); err != nil {
			return nil, err
		}

		return rep, nil
	}

	rep = report.Downloaded(j, int64(len(bs)))

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

func (cli *Client) useCache(rep *report.Report) bool {
	if rep == nil {
		return false
	}
	return !shouldRetry(cli.retry, rep)
}

func shouldRetry(cfg *retry.RetryConfig, rep *report.Report) bool {
	if cfg == nil || cfg.ShouldRetry == nil {
		return retry.DefaultRetryConfig.ShouldRetry(rep)
	}
	return cfg.ShouldRetry(rep)
}
