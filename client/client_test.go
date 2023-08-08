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
	"github.com/stretchr/testify/require"

	"github.com/hiroara/drawin/client"
	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/reporter"
)

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
