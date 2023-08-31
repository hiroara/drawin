package drawin

import (
	"github.com/hiroara/drawin/job"
)

type Output interface {
	Add(*Report, []byte) error
	Get(*job.Job) (*Report, error)
	Initialize() error
}
