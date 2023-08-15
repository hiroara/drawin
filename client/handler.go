package client

import (
	"context"
	"net/http"

	"github.com/hiroara/drawin/job"
)

type Handler interface {
	Match(*job.Job) bool
	ShouldRetry(error) bool
	Get(context.Context, *job.Job) ([]byte, error)
}

var DefaultHandlers = []Handler{&HTTPHandler{client: http.DefaultClient}}
