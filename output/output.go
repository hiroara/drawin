package output

import (
	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/report"
)

type Output interface {
	Add(*report.Report, []byte) error
	Get(*job.Job) (*report.Report, error)
	Initialize() error
}
