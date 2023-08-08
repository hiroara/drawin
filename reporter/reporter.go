package reporter

import (
	"io"
)

type Reporter interface {
	io.Closer
	Write(*Report) error
}
