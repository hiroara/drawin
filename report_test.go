package drawin_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hiroara/drawin"
	"github.com/hiroara/drawin/job"
)

func TestDownloadedReport(t *testing.T) {
	t.Parallel()

	j := &job.Job{Name: "image1.jpg"}

	dr := drawin.DownloadedReport(j, 256)
	if assert.NotNil(t, dr) {
		assert.Equal(t, *j, dr.Job)
		assert.Equal(t, drawin.DownloadedResult, dr.Result)
		assert.Equal(t, int64(256), dr.ContentLength)
		assert.Empty(t, dr.Failure)
	}
}

func TestCachedReport(t *testing.T) {
	t.Parallel()

	j := &job.Job{Name: "image1.jpg"}

	dr := drawin.CachedReport(j, 512)
	if assert.NotNil(t, dr) {
		assert.Equal(t, *j, dr.Job)
		assert.Equal(t, drawin.CachedResult, dr.Result)
		assert.Equal(t, int64(512), dr.ContentLength)
		assert.Empty(t, dr.Failure)
	}
}

func TestFailedReport(t *testing.T) {
	t.Parallel()

	j := &job.Job{Name: "image1.jpg"}
	err := errors.New("test error")

	dr := drawin.FailedReport(j, err, true)
	if assert.NotNil(t, dr) {
		assert.Equal(t, *j, dr.Job)
		assert.Equal(t, drawin.FailedResult, dr.Result)
		assert.Empty(t, dr.ContentLength)
		if assert.NotEmpty(t, dr.Failure) {
			assert.Equal(t, err.Error(), dr.Failure.Error)
			assert.True(t, dr.Failure.Permanent)
		}
	}
}
