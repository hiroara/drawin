package output_test

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/output"
	"github.com/hiroara/drawin/report"
)

func TestDirectoryOutput(t *testing.T) {
	t.Parallel()

	j := &job.Job{Name: "file1.txt"}
	data := []byte("test value")

	buildOutput := func(t *testing.T) (string, *output.DirectoryOutput) {
		dir := filepath.Join(t.TempDir(), "out")
		out := output.NewDirectory(dir)
		return dir, out
	}

	t.Run("NoDataCase", func(t *testing.T) {
		t.Parallel()

		_, out := buildOutput(t)
		require.NoError(t, out.Initialize())
		rep, err := out.Get(j)
		require.NoError(t, err)
		assert.Nil(t, rep)
	})

	t.Run("DownloadedReportCase", func(t *testing.T) {
		t.Parallel()

		dir, out := buildOutput(t)
		require.NoError(t, out.Initialize())

		require.NoError(t, out.Add(report.Downloaded(j, int64(len(data))), data))

		f, err := os.Open(filepath.Join(dir, "file1.txt"))
		assert.NoError(t, err)
		bs, err := io.ReadAll(f)
		assert.NoError(t, err)
		assert.Equal(t, data, bs)

		rep, err := out.Get(j)
		require.NoError(t, err)
		if assert.NotNil(t, rep) {
			assert.Equal(t, report.CachedResult, rep.Result)
			assert.Equal(t, int64(len(data)), rep.ContentLength)
		}
	})

	t.Run("FailedReportCase", func(t *testing.T) {
		t.Parallel()

		dir, out := buildOutput(t)
		require.NoError(t, out.Initialize())

		require.NoError(t, out.Add(report.Failed(j, errors.New("test error"), true), nil))
		_, err := os.Stat(filepath.Join(dir, "file1.txt"))
		require.ErrorIs(t, err, os.ErrNotExist)
	})
}
