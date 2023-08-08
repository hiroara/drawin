package reporter_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/reporter"
)

func TestJSONReporter(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "report.jsonl")
	rep, err := reporter.OpenJSON(path)
	require.NoError(t, err)

	rep.Write(reporter.DownloadedReport(&job.Job{Name: "image1.jpg"}, 256))
	rep.Write(reporter.CachedReport(&job.Job{Name: "image2.jpg"}))

	require.NoError(t, rep.Close())

	f, err := os.Open(path)
	require.NoError(t, err)
	bs, err := io.ReadAll(f)
	require.NoError(t, err)

	assert.Equal(t, "{\"job\":{\"url\":\"\",\"name\":\"image1.jpg\"},\"result\":\"downloaded\",\"contentLength\":256}\n{\"job\":{\"url\":\"\",\"name\":\"image2.jpg\"},\"result\":\"cached\"}\n", string(bs))
}
