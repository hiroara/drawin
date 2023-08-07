package reporter

import (
	"io"

	"github.com/hiroara/drawin/job"
)

type Reporter interface {
	io.Closer
	Write(*job.Job) error
}
