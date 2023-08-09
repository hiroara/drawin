package client_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/drawin/client"
	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/reporter"
)

type customHandler struct {
	mock.Mock
}

func (h *customHandler) Match(j *job.Job) bool {
	args := h.Called(j)
	return args.Bool(0)
}

func (h *customHandler) Get(ctx context.Context, j *job.Job) ([]byte, error) {
	args := h.Called(ctx, j)
	return args.Get(0).([]byte), args.Error(1)
}

func TestDownload(t *testing.T) {
	t.Run("ResponseStatusCode=OK", func(t *testing.T) {
		dirpath := filepath.Join(t.TempDir(), "test-out")
		dir := client.NewDirectory(dirpath)
		cli, err := client.Build(dir)
		require.NoError(t, err)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprintln(w, "Successful")
		}))
		defer srv.Close()

		j := &job.Job{Name: "image1.jpg", URL: srv.URL}
		rep, err := cli.Download(context.Background(), j)
		require.NoError(t, err)
		assert.Equal(t, *j, rep.Job)
		assert.Equal(t, reporter.Downloaded, rep.Result)

		f, err := os.Open(filepath.Join(dirpath, "image1.jpg"))
		require.NoError(t, err)
		buf := bytes.NewBuffer(nil)
		_, err = io.Copy(buf, f)
		require.NoError(t, err)
		assert.Equal(t, "Successful\n", buf.String())
		assert.Equal(t, int64(11), rep.ContentLength)

		rep, err = cli.Download(context.Background(), j)
		require.NoError(t, err)
		assert.Equal(t, *j, rep.Job)
		assert.Equal(t, reporter.Cached, rep.Result)
	})

	t.Run("ResponseStatusCode=NotFound", func(t *testing.T) {
		dirpath := filepath.Join(t.TempDir(), "test-out")
		dir := client.NewDirectory(dirpath)
		cli, err := client.Build(dir)
		require.NoError(t, err)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(404)
			fmt.Fprintln(w, "Not Found")
		}))
		defer srv.Close()

		j := &job.Job{Name: "image1.jpg", URL: srv.URL}
		rep, err := cli.Download(context.Background(), j)
		require.NoError(t, err)
		assert.Equal(t, reporter.Failed, rep.Result)

		_, err = os.Stat(filepath.Join(dirpath, "image1.jpg"))
		require.ErrorIs(t, err, os.ErrNotExist)
	})
}

func TestDownloadWithCustomHandler(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	dirpath := filepath.Join(t.TempDir(), "test-out")
	dir := client.NewDirectory(dirpath)

	h := &customHandler{}
	cli, err := client.Build(dir, client.WithHandlers(h))
	require.NoError(t, err)

	j := &job.Job{Name: "test1.jpg", URL: "https://example.com/test1.jpg"}
	data := []byte("downloaded content")

	h.On("Match", j).Return(true).Once()
	h.On("Get", ctx, j).Return(data, nil).Once()

	rep, err := cli.Download(ctx, j)
	require.NoError(t, err)
	if assert.NotNil(t, rep) {
		assert.Equal(t, reporter.Downloaded, rep.Result)
		assert.Equal(t, int64(len(data)), rep.ContentLength)
	}
	h.AssertExpectations(t)
}
