package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"

	"github.com/hiroara/drawin/database"
	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/reporter"
	"github.com/hiroara/drawin/store"
)

var errNoMatchingReport = errors.New("report not found")
var errNoSuccessfulReport = errors.New("report found but its content has not been downloaded")
var errNoMatchingData = errors.New("data not found")

func read(path string, urls []string) (*flow.Flow, error) {
	db, err := database.Open(path, &database.Options{Create: false})
	if err != nil {
		return nil, err
	}
	s := store.New(db)

	src := source.FromSlice(urls)

	rs := task.Connect(
		src.AsTask(),
		pipe.Map(func(ctx context.Context, url string) (*reporter.Report, error) {
			rep, err := s.Get(&job.Job{URL: url})
			if err != nil {
				return nil, err
			}
			if rep == nil {
				return nil, fmt.Errorf("%w with URL: %s", errNoMatchingReport, url)
			}
			if rep.Result != reporter.Downloaded {
				return nil, fmt.Errorf("%w with URL: %s", errNoSuccessfulReport, url)
			}
			return rep, nil
		}).AsTask(),
		0,
	)

	ds := task.Connect(
		rs,
		pipe.Map(func(ctx context.Context, rep *reporter.Report) ([]byte, error) {
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

	return flow.FromTask(sin), nil
}
