package reader

import (
	"bufio"
	"context"
	"io"
	"os"

	"github.com/hiroara/carbo/task"
)

func Read(ctx context.Context, path string, out chan<- string) error {
	r, err := open(path)
	if err != nil {
		return err
	}
	defer r.Close()

	return emitLines(ctx, r, out)
}

func open(path string) (io.ReadCloser, error) {
	if path == "-" {
		return os.Stdin, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func emitLines(ctx context.Context, r io.Reader, out chan<- string) error {
	sc := bufio.NewScanner(r)

	for sc.Scan() {
		if err := task.Emit(ctx, out, sc.Text()); err != nil {
			return err
		}
	}
	return nil
}
