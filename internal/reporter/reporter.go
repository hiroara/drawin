package reporter

import (
	"io"

	"github.com/hiroara/drawin"
)

type Reporter interface {
	io.Closer
	Write(*drawin.Report) error
}
