package job

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

type namer struct {
	url      *url.URL
	attempts int
}

func newNamer(j *Job) (*namer, error) {
	u, err := url.ParseRequestURI(j.URL)
	if err != nil {
		return nil, err
	}
	return &namer{url: u}, nil
}

func (n *namer) name() string {
	bs := path.Base(n.url.Path)
	if n.attempts == 0 {
		return bs
	}
	ts := strings.Split(bs, ".")
	if len(ts) == 1 {
		return fmt.Sprintf("%s.%d", ts[0], n.attempts)
	}
	fn := strings.Join(ts[:len(ts)-1], ".")
	return fmt.Sprintf("%s.%d.%s", fn, n.attempts, ts[len(ts)-1])
}

func (n *namer) next() {
	n.attempts++
}
