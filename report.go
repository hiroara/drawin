package main

import (
	"context"
	"encoding/json"
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

func getReports(s *store.Store, urls []string) source.Source[*reporter.Report] {
	src := source.FromSlice(urls)

	return task.Connect(
		src.AsTask(),
		pipe.Map(func(ctx context.Context, url string) (*reporter.Report, error) {
			rep, err := s.Get(&job.Job{URL: url})
			if err != nil {
				return nil, err
			}
			if rep == nil {
				return nil, fmt.Errorf("%w with URL: %s", errNoMatchingReport, url)
			}
			return rep, nil
		}).AsTask(),
		0,
	)
}

func report(path string, urls []string) (*flow.Flow, error) {
	s, err := store.Open(path, &database.Options{Create: false})
	if err != nil {
		return nil, err
	}

	reps := getReports(s, urls)

	enc := json.NewEncoder(os.Stdout)
	sin := task.Connect(
		reps.AsTask(),
		sink.ElementWise(func(ctx context.Context, rep *reporter.Report) error {
			return enc.Encode(rep)
		}).AsTask(),
		0,
	)

	return flow.FromTask(sin), nil
}
