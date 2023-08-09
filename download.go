package main

import (
	"context"
	"path/filepath"

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

func download(paths []string, outdir, reportPath string, useStore bool, concurrency int) (*flow.Flow, error) {
	var out client.Output
	if useStore {
		db, err := database.Open(filepath.Join(outdir, "drawin.db"))
		if err != nil {
			return nil, err
		}
		defer db.Close()

		out = store.New(db)
	} else {
		out = client.NewDirectory(outdir)
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

	sin := task.Connect(
		reps,
		sink.ElementWise(func(ctx context.Context, rep *reporter.Report) error { return repr.Write(rep) }).AsTask(),
		0,
	)

	return flow.FromTask(sin), nil
}
