package main

import (
	"context"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"
	"golang.org/x/sync/errgroup"

	"github.com/hiroara/drawin/client"
	"github.com/hiroara/drawin/database"
	"github.com/hiroara/drawin/downloader"
	"github.com/hiroara/drawin/reader"
	"github.com/hiroara/drawin/reporter"
	"github.com/hiroara/drawin/store"
)

func download(paths []string, outStr, reportPath string, concurrency int) (*flow.Flow, error) {
	o, err := parseOutput(outStr)
	if err != nil {
		return nil, err
	}

	var out client.Output
	closeOut := func() {}
	switch o.typ {
	case storeType:
		db, err := database.Open(o.path)
		if err != nil {
			return nil, err
		}

		out = store.New(db)
		closeOut = func() { db.Close() }
	case directoryType:
		out = client.NewDirectory(o.path)
	}

	cli, err := client.Build(out)
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

	reps := task.Connect(
		urls,
		pipe.FromFn(func(ctx context.Context, in <-chan string, out chan<- *reporter.Report) error {
			m := make(chan *reporter.Report)

			d, err := downloader.New(cli, m)
			if err != nil {
				return err
			}

			grp, ctx := errgroup.WithContext(ctx)

			grp.Go(func() error { return d.Run(ctx, in) })
			grp.Go(func() error {
				for rep := range m {
					if err := task.Emit(ctx, out, rep); err != nil {
						return err
					}
				}
				return nil
			})

			return grp.Wait()
		}).AsTask(),
		32,
	)
	reps.Defer(closeOut)

	sin := task.Connect(
		reps,
		sink.ElementWise(func(ctx context.Context, rep *reporter.Report) error { return repr.Write(rep) }).AsTask(),
		0,
	)

	return flow.FromTask(sin), nil
}
