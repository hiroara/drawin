package reporter_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/drawin/downloader/report"
	"github.com/hiroara/drawin/internal/reporter"
	"github.com/hiroara/drawin/job"
)

func TestJSONReporter(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "report.jsonl")
	rep, err := reporter.OpenJSON(path)
	require.NoError(t, err)

	rep.Write(report.Downloaded(&job.Job{Name: "image1.jpg"}, 256))
	rep.Write(report.Cached(&job.Job{Name: "image2.jpg"}, 512))

	require.NoError(t, rep.Close())

	f, err := os.Open(path)
	require.NoError(t, err)
	bs, err := io.ReadAll(f)
	require.NoError(t, err)

	assert.Equal(t, "{\"job\":{\"url\":\"\",\"name\":\"image1.jpg\"},\"result\":\"downloaded\",\"contentLength\":256}\n{\"job\":{\"url\":\"\",\"name\":\"image2.jpg\"},\"result\":\"cached\",\"contentLength\":512}\n", string(bs))
}
