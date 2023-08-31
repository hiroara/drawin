package handler

import (
	"context"
	"net/http"

	httphandler "github.com/hiroara/drawin/handler/http"
	"github.com/hiroara/drawin/job"
)

type Handler interface {
	Match(*job.Job) bool
	ShouldRetry(error) bool
	Get(context.Context, *job.Job) ([]byte, error)
}

var DefaultHandlers = []Handler{httphandler.New(http.DefaultClient)}
