package client_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/drawin"
	"github.com/hiroara/drawin/client"
	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/store"
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

func TestDownload(t *testing.T) {
	t.Run("ResponseStatusCode=OK", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "test.db")
		s, err := store.Open(path, nil)
		require.NoError(t, err)
		cli, err := client.Build(s, nil, nil)
		require.NoError(t, err)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprintln(w, "Successful")
		}))
		defer srv.Close()

		j := &job.Job{Name: "image1.jpg", URL: srv.URL}
		rep, err := cli.Download(context.Background(), j)
		require.NoError(t, err)
		assert.Equal(t, *j, rep.Job)
		assert.Equal(t, drawin.DownloadedResult, rep.Result)

		rep, err = s.Get(j)
		require.NoError(t, err)
		if assert.NotNil(t, rep) {
			assert.Equal(t, drawin.DownloadedResult, rep.Result)
		}

		data, err := s.Read(rep)
		require.NoError(t, err)

		assert.Equal(t, "Successful\n", string(data))
		assert.Equal(t, int64(11), rep.ContentLength)

		rep, err = cli.Download(context.Background(), j)
		require.NoError(t, err)
		assert.Equal(t, *j, rep.Job)
		assert.Equal(t, drawin.CachedResult, rep.Result)
	})

	t.Run("ResponseStatusCode=NotFound", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "test.db")
		s, err := store.Open(path, nil)
		require.NoError(t, err)
		cli, err := client.Build(s, nil, nil)
		require.NoError(t, err)

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(404)
			fmt.Fprintln(w, "Not Found")
		}))
		defer srv.Close()

		j := &job.Job{Name: "image1.jpg", URL: srv.URL}
		rep, err := cli.Download(context.Background(), j)
		require.NoError(t, err)
		assert.Equal(t, drawin.FailedResult, rep.Result)

		rep, err = s.Get(j)
		require.NoError(t, err)
		if assert.NotNil(t, rep) {
			assert.Equal(t, drawin.FailedResult, rep.Result)
		}

		data, err := s.Read(rep)
		require.NoError(t, err)
		assert.Nil(t, data)
	})
}
