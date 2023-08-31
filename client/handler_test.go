package client_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/drawin"
	"github.com/hiroara/drawin/client"
	"github.com/hiroara/drawin/job"
)

func TestCustomHandler(t *testing.T) {
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
		assert.Equal(t, drawin.DownloadedResult, rep.Result)
		assert.Equal(t, int64(len(data)), rep.ContentLength)
	}
	h.AssertExpectations(t)
}
