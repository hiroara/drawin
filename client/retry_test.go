package client_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/drawin/client"
	"github.com/hiroara/drawin/downloader"
	"github.com/hiroara/drawin/downloader/report"
	"github.com/hiroara/drawin/job"
)

func buildClient(t *testing.T, cfg *client.RetryConfig) (*client.Client, error) {
	dirpath := filepath.Join(t.TempDir(), "test-out")
	dir := client.NewDirectory(dirpath)

	return client.Build(dir, client.WithRetryConfig(cfg))
}

func TestDefaultRetryConfig(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, "Successful")
	}))
	defer srv.Close()

	j := &job.Job{Name: "image1.jpg", URL: srv.URL}

	cli, err := buildClient(t, nil)
	require.NoError(t, err)

	rep, err := cli.Download(context.Background(), j)
	require.NoError(t, err)
	if assert.NotNil(t, rep) {
		assert.Equal(t, *j, rep.Job)
		assert.Equal(t, report.DownloadedResult, rep.Result)
	}

	rep, err = cli.Download(context.Background(), j)
	require.NoError(t, err)
	if assert.NotNil(t, rep) {
		assert.Equal(t, *j, rep.Job)
		assert.Equal(t, report.CachedResult, rep.Result)
	}
}

func TestCustomRetryConfig(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, "Successful")
	}))
	defer srv.Close()

	j := &job.Job{Name: "image1.jpg", URL: srv.URL}
	cli, err := buildClient(t, &client.RetryConfig{
		// Always retry
		ShouldRetry: func(rep *downloader.Report) bool { return true },
	})
	require.NoError(t, err)

	rep, err := cli.Download(context.Background(), j)
	require.NoError(t, err)
	if assert.NotNil(t, rep) {
		assert.Equal(t, *j, rep.Job)
		assert.Equal(t, report.DownloadedResult, rep.Result)
	}

	rep, err = cli.Download(context.Background(), j)
	require.NoError(t, err)
	if assert.NotNil(t, rep) {
		assert.Equal(t, *j, rep.Job)
		assert.Equal(t, report.DownloadedResult, rep.Result) // Download again
	}
}
