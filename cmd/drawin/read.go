package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/task"
	"github.com/hiroara/drawin"

	"github.com/hiroara/drawin/database"
	"github.com/hiroara/drawin/store"
)

var errNoMatchingData = errors.New("data not found")

func runRead(path string, urls []string) (*flow.Flow, error) {
	s, err := store.Open(path, &database.Options{Create: false})
	if err != nil {
		return nil, err
	}

	reps := getReports(s, urls)

	ds := task.Connect(
		reps.AsTask(),
		pipe.Map(func(ctx context.Context, rep *drawin.Report) ([]byte, error) {
			if rep.Result != drawin.DownloadedResult {
				return nil, fmt.Errorf("%w with URL: %s", errNoSuccessfulReport, rep.URL)
			}

			bs, err := s.Read(rep)
			if err != nil {
				return nil, err
			}
			if bs == nil {
				return nil, fmt.Errorf("%w with URL: %s", errNoMatchingData, rep.URL)
			}
			return bs, nil
		}).AsTask(),
		0,
	)

	sin := task.Connect(
		ds,
		sink.ElementWise(func(ctx context.Context, bs []byte) error {
			_, err := os.Stdout.Write(bs)
			return err
		}).AsTask(),
		0,
	)
	sin.Defer(func() { s.Close() })

	return flow.FromTask(sin), nil
}
