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
	path := filepath.Join(t.TempDir(), "report.jsonl")
	rep, err := reporter.OpenJSON(path)
	require.NoError(t, err)

	rep.Write(&job.Job{Name: "image1.jpg"})
	rep.Write(&job.Job{Name: "image2.jpg"})

	require.NoError(t, rep.Close())

	f, err := os.Open(path)
	require.NoError(t, err)
	bs, err := io.ReadAll(f)
	require.NoError(t, err)

	assert.Equal(t, "{\"url\":\"\",\"name\":\"image1.jpg\",\"contentLength\":0,\"action\":\"\"}\n{\"url\":\"\",\"name\":\"image2.jpg\",\"contentLength\":0,\"action\":\"\"}\n", string(bs))
}
