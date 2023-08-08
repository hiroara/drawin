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
)

func TestDirectoryOutput(t *testing.T) {
	t.Parallel()

	dir := filepath.Join(t.TempDir(), "out")
	out := client.NewDirectory(dir)
	require.NoError(t, out.Prepare())

	j := &job.Job{Name: "file1.txt"}

	ok, err := out.Check(j)
	require.NoError(t, err)
	assert.False(t, ok)

	require.NoError(t, out.Add(j, []byte("test value")))

	f, err := os.Open(filepath.Join(dir, "file1.txt"))
	assert.NoError(t, err)
	bs, err := io.ReadAll(f)
	assert.NoError(t, err)
	assert.Equal(t, []byte("test value"), bs)

	ok, err = out.Check(j)
	require.NoError(t, err)
	assert.True(t, ok)
}
