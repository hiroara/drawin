package handler_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/drawin"
	"github.com/hiroara/drawin/client"
	"github.com/hiroara/drawin/handler"
	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/output"
)

type customHandler struct {
	mock.Mock
}

func (h *customHandler) Match(j *job.Job) bool {
	args := h.Called(j)
	return args.Bool(0)
}

func (h *customHandler) ShouldRetry(err error) bool {
	args := h.Called(err)
	return args.Bool(0)
}

func (h *customHandler) Get(ctx context.Context, j *job.Job) ([]byte, error) {
	args := h.Called(ctx, j)
	return args.Get(0).([]byte), args.Error(1)
}

func TestCustomHandler(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dirpath := filepath.Join(t.TempDir(), "test-out")
	dir := output.NewDirectory(dirpath)

	h := &customHandler{}
	cli, err := client.Build(dir, []handler.Handler{h}, nil)
	require.NoError(t, err)

	j := &job.Job{Name: "test1.jpg", URL: "https://example.com/test1.jpg"}
	data := []byte("downloaded content")

	h.On("Match", j).Return(true).Once()
	h.On("Get", ctx, j).Return(data, nil).Once()

	rep, err := cli.Download(ctx, j)
	require.NoError(t, err)
	if assert.NotNil(t, rep) {
		assert.Equal(t, drawin.DownloadedResult, rep.Result)
		assert.Equal(t, int64(len(data)), rep.ContentLength)
	}
	h.AssertExpectations(t)
}
