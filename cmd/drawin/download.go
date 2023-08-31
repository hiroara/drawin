package main

import (
	"context"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"

	"github.com/hiroara/drawin"
	"github.com/hiroara/drawin/client"
	"github.com/hiroara/drawin/internal/reporter"
	"github.com/hiroara/drawin/reader"
	"github.com/hiroara/drawin/store"
)

func runDownload(paths []string, outStr, reportPath string, concurrency int) (*flow.Flow, error) {
	o, err := parseOutput(outStr)
	if err != nil {
		return nil, err
	}

	var out client.Output
	closeOut := func() {}
	switch o.typ {
	case storeType:
		s, err := store.Open(o.path, nil)
		if err != nil {
			return nil, err
		}

		out = s
		closeOut = func() { s.Close() }
	case directoryType:
		out = client.NewDirectory(o.path)
	}

	cli, err := client.Build(out, nil, nil)
	if err != nil {
		return nil, err
	}

	repr, err := reporter.OpenJSON(reportPath)
	if err != nil {
		return nil, err
	}

	src := source.FromSlice(paths)

	urls := task.Connect(
		src.AsTask(),
		pipe.FromFn(func(ctx context.Context, in <-chan string, out chan<- string) error {
			for path := range in {
				if err := reader.Read(ctx, path, out); err != nil {
					return err
				}
			}
			return nil
		}).AsTask(),
		0,
	)

	d, err := drawin.NewDownloader(cli)
	if err != nil {
		return nil, err
	}

	reps := task.Connect(
		urls,
		d.AsTask(),
		32,
	)
	reps.Defer(closeOut)
	reps.Defer(func() { d.Close() })

	sin := task.Connect(
		reps,
		sink.ElementWise(func(ctx context.Context, rep *drawin.Report) error { return repr.Write(rep) }).AsTask(),
		0,
	)

	return flow.FromTask(sin), nil
}
