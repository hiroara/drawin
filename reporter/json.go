package reporter

import (
	"encoding/json"
	"io"
	"os"

	"github.com/hiroara/drawin/job"
)

type jsonReporter struct {
	out           io.WriteCloser
	writer        *json.Encoder
	headerWritten bool
}

var headers = []string{"url", "path", "action"}

func NewJSON(out io.WriteCloser) Reporter {
	return &jsonReporter{out: out, writer: json.NewEncoder(out), headerWritten: false}
}

func OpenJSON(path string) (Reporter, error) {
	if path == "-" {
		return NewJSON(os.Stdout), nil
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_EXCL|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return NewJSON(f), nil
}

func (r *jsonReporter) Write(j *job.Job) error {
	return r.write(j)
}

func (r *jsonReporter) write(j *job.Job) error {
	return r.writer.Encode(j)
}

func (r *jsonReporter) Close() error {
	return r.out.Close()
}