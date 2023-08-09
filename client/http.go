package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hiroara/drawin/job"
)

type HTTPHandler struct {
	client *http.Client
}

func NewHTTPHandler(cli *http.Client) *HTTPHandler {
	return &HTTPHandler{client: cli}
}

func (h *HTTPHandler) Match(j *job.Job) bool {
	return strings.HasPrefix(j.URL, "http://") || strings.HasPrefix(j.URL, "https://")
}

func (h *HTTPHandler) Get(ctx context.Context, j *job.Job) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", j.URL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 || resp.StatusCode < 200 {
		return nil, fmt.Errorf("%w: %d", ErrUnexpectedResponseStatus, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
