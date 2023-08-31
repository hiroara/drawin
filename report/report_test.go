package report_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/report"
)

func TestDownloaded(t *testing.T) {
	t.Parallel()

	j := &job.Job{Name: "image1.jpg"}

	dr := report.Downloaded(j, 256)
	if assert.NotNil(t, dr) {
		assert.Equal(t, *j, dr.Job)
		assert.Equal(t, report.DownloadedResult, dr.Result)
		assert.Equal(t, int64(256), dr.ContentLength)
		assert.Empty(t, dr.Failure)
	}
}

func TestCached(t *testing.T) {
	t.Parallel()

	j := &job.Job{Name: "image1.jpg"}

	dr := report.Cached(j, 512)
	if assert.NotNil(t, dr) {
		assert.Equal(t, *j, dr.Job)
		assert.Equal(t, report.CachedResult, dr.Result)
		assert.Equal(t, int64(512), dr.ContentLength)
		assert.Empty(t, dr.Failure)
	}
}

func TestFailed(t *testing.T) {
	t.Parallel()

	j := &job.Job{Name: "image1.jpg"}
	err := errors.New("test error")

	dr := report.Failed(j, err, true)
	if assert.NotNil(t, dr) {
		assert.Equal(t, *j, dr.Job)
		assert.Equal(t, report.FailedResult, dr.Result)
		assert.Empty(t, dr.ContentLength)
		if assert.NotEmpty(t, dr.Failure) {
			assert.Equal(t, err.Error(), dr.Failure.Error)
			assert.True(t, dr.Failure.Permanent)
		}
	}
}
