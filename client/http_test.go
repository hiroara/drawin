package client_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/drawin/client"
	"github.com/hiroara/drawin/job"
)

func TestHTTPHandlerMatch(t *testing.T) {
	t.Parallel()

	cli := client.NewHTTPHandler(http.DefaultClient)

	assert.True(t, cli.Match(&job.Job{URL: "http://example.com/test1.txt"}))
	assert.True(t, cli.Match(&job.Job{URL: "https://example.com/test1.txt"}))
	assert.False(t, cli.Match(&job.Job{URL: "ftp://example.com/test1.txt"}))
	assert.False(t, cli.Match(&job.Job{URL: "file:///etc/hosts"}))
}

func TestHTTPHandlerGet(t *testing.T) {
	t.Parallel()

	t.Run("ResponseStatusCode=Successful", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprintln(w, "Successful")
		}))
		defer srv.Close()

		j := &job.Job{Name: "image1.jpg", URL: srv.URL}

		cli := client.NewHTTPHandler(http.DefaultClient)

		bs, err := cli.Get(context.Background(), j)
		require.NoError(t, err)
		assert.Len(t, bs, 11)
	})

	t.Run("ResponseStatusCode=Redirection", func(t *testing.T) {
		t.Parallel()

		srv1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprintln(w, "Successful")
		}))
		defer srv1.Close()

		srv2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			http.Redirect(w, req, srv1.URL, 302)
		}))
		defer srv2.Close()

		j := &job.Job{Name: "image1.jpg", URL: srv2.URL}

		cli := client.NewHTTPHandler(http.DefaultClient)

		bs, err := cli.Get(context.Background(), j)
		require.NoError(t, err)
		assert.Len(t, bs, 11)
	})

	t.Run("ResponseStatusCode=ClientError", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(400)
			fmt.Fprintln(w, "Client error")
		}))
		defer srv.Close()

		j := &job.Job{Name: "image1.jpg", URL: srv.URL}

		cli := client.NewHTTPHandler(http.DefaultClient)

		_, err := cli.Get(context.Background(), j)
		require.ErrorIs(t, err, client.ErrClientError)
	})

	t.Run("ResponseStatusCode=Unexpected", func(t *testing.T) {
		t.Parallel()

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(500)
			fmt.Fprintln(w, "Server internal error")
		}))
		defer srv.Close()

		j := &job.Job{Name: "image1.jpg", URL: srv.URL}

		cli := client.NewHTTPHandler(http.DefaultClient)

		_, err := cli.Get(context.Background(), j)
		require.ErrorIs(t, err, client.ErrUnexpectedResponseStatus)
	})
}
