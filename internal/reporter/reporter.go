package reporter

import (
	"io"

	"github.com/hiroara/drawin/downloader/report"
)

type Reporter interface {
	io.Closer
	Write(*report.Report) error
}
