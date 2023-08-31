package http

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hiroara/drawin/job"
)

type Handler struct {
	client *http.Client
}

func New(cli *http.Client) *Handler {
	return &Handler{client: cli}
}

func (h *Handler) Match(j *job.Job) bool {
	return strings.HasPrefix(j.URL, "http://") || strings.HasPrefix(j.URL, "https://")
}

func (h *Handler) ShouldRetry(err error) bool {
	return !errors.Is(err, ErrClientError)
}

var ErrUnexpectedResponseStatus = errors.New("received unexpected HTTP response status code")
var ErrClientError = errors.New("received client error response")

func (h *Handler) Get(ctx context.Context, j *job.Job) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", j.URL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	st := resp.StatusCode
	switch {
	case st >= 200 && st < 300:
	case st >= 400 && st < 500:
		return nil, fmt.Errorf("%w: %d", ErrClientError, resp.StatusCode)
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedResponseStatus, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
