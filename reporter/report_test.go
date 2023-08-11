package reporter_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/reporter"
)

func TestDownloadedReport(t *testing.T) {
	t.Parallel()

	j := &job.Job{Name: "image1.jpg"}

	dr := reporter.DownloadedReport(j, 256)
	if assert.NotNil(t, dr) {
		assert.Equal(t, *j, dr.Job)
		assert.Equal(t, reporter.Downloaded, dr.Result)
		assert.Equal(t, int64(256), dr.ContentLength)
		assert.Empty(t, dr.Failure)
	}
}

func TestCachedReport(t *testing.T) {
	t.Parallel()

	j := &job.Job{Name: "image1.jpg"}

	dr := reporter.CachedReport(j, 512)
	if assert.NotNil(t, dr) {
		assert.Equal(t, *j, dr.Job)
		assert.Equal(t, reporter.Cached, dr.Result)
		assert.Equal(t, int64(512), dr.ContentLength)
		assert.Empty(t, dr.Failure)
	}
}

func TestFailedReport(t *testing.T) {
	t.Parallel()

	j := &job.Job{Name: "image1.jpg"}
	err := errors.New("test error")

	dr := reporter.FailedReport(j, err, true)
	if assert.NotNil(t, dr) {
		assert.Equal(t, *j, dr.Job)
		assert.Equal(t, reporter.Failed, dr.Result)
		assert.Empty(t, dr.ContentLength)
		if assert.NotEmpty(t, dr.Failure) {
			assert.Equal(t, err.Error(), dr.Failure.Error)
			assert.True(t, dr.Failure.Permanent)
		}
	}
}
