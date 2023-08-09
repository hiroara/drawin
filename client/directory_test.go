package client_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/drawin/client"
	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/reporter"
)

func TestDirectoryOutput(t *testing.T) {
	t.Parallel()

	dir := filepath.Join(t.TempDir(), "out")
	out := client.NewDirectory(dir)
	require.NoError(t, out.Initialize())

	j := &job.Job{Name: "file1.txt"}
	data := []byte("test value")

	rep, err := out.Get(j)
	require.NoError(t, err)
	assert.Nil(t, rep)

	require.NoError(t, out.Add(reporter.DownloadedReport(j, int64(len(data))), data))

	f, err := os.Open(filepath.Join(dir, "file1.txt"))
	assert.NoError(t, err)
	bs, err := io.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, data, bs)

	rep, err = out.Get(j)
	require.NoError(t, err)
	if assert.NotNil(t, rep) {
		assert.Equal(t, reporter.Cached, rep.Result)
		assert.Equal(t, int64(len(data)), rep.ContentLength)
	}
}
