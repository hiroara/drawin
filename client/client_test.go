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
)

func TestCreateDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "test-out")
	cli := client.New(dir)

	_, err := os.Stat(dir)
	require.ErrorIs(t, err, os.ErrNotExist)

	require.NoError(t, cli.CreateDir())

	stat, err := os.Stat(dir)
	require.NoError(t, err)
	if assert.NotNil(t, stat) {
		assert.True(t, stat.IsDir())
	}
}

func TestDownload(t *testing.T) {
	t.Run("ResponseStatusCode=OK", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "test-out")
		cli := client.New(dir)
		require.NoError(t, cli.CreateDir())

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprintln(w, "Successful")
		}))
		defer srv.Close()

		j := &job.Job{Name: "image1.jpg", URL: srv.URL}
		err := cli.Download(context.Background(), j)
		require.NoError(t, err)
		assert.Equal(t, job.DownloadAction, j.Action)

		f, err := os.Open(filepath.Join(dir, "image1.jpg"))
		require.NoError(t, err)
		buf := bytes.NewBuffer(nil)
		_, err = io.Copy(buf, f)
		require.NoError(t, err)
		assert.Equal(t, "Successful\n", buf.String())
		assert.Equal(t, int64(11), j.ContentLength)
	})

	t.Run("ResponseStatusCode=NotFound", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "test-out")
		cli := client.New(dir)
		require.NoError(t, cli.CreateDir())

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(404)
			fmt.Fprintln(w, "Not Found")
		}))
		defer srv.Close()

		j := &job.Job{Name: "image1.jpg", URL: srv.URL}
		err := cli.Download(context.Background(), j)
		require.Error(t, err)
		assert.Empty(t, j.Action)
		assert.Empty(t, j.ContentLength)

		_, err = os.Stat(filepath.Join(dir, "image1.jpg"))
		require.ErrorIs(t, err, os.ErrNotExist)
	})
}
