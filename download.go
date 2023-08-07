package main

import (
	"context"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"

	"github.com/hiroara/drawin/client"
	"github.com/hiroara/drawin/database"
	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/reader"
	"github.com/hiroara/drawin/reporter"
)

func downloadFiles(paths []string, reportPath string, cacheDB *database.DB) (*flow.Flow, error) {
	cli := client.New(outdir)
	if err := cli.CreateDir(); err != nil {
		return nil, err
	}

	ps := source.FromSlice(paths)
	ls := task.Connect(
		ps.AsTask(),
		pipe.FromFn(reader.Read).AsTask(),
		0,
	)
	lbs := task.Connect(
		ls.AsTask(),
		pipe.Batch[string](32).AsTask(),
		0,
	)
	jbs := task.Connect(
		lbs.AsTask(),
		pipe.Map(func(ctx context.Context, urls []string) ([]*job.Job, error) {
			jobs := make([]*job.Job, 0, len(urls))
			err := database.Update(cacheDB, func(buc *database.Bucket[*job.Job]) error {
				s := job.NewStore(buc)
				for _, u := range urls {
					j, err := s.CreateJob(u)
					if err != nil {
						return err
					}
					if j != nil {
						jobs = append(jobs, j)
					}
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
			return jobs, nil
		}).AsTask(),
		1,
	)
	js := task.Connect(
		jbs.AsTask(),
		pipe.FlattenSlice[*job.Job]().AsTask(),
		0,
	)
	js = task.Connect(
		js.AsTask(),
		pipe.Tap(cli.Download).Concurrent(concurrency).AsTask(),
		0,
	)
	js = task.Connect(
		js.AsTask(),
		pipe.Tap(func(ctx context.Context, j *job.Job) error {
			return database.Update(cacheDB, func(buc *database.Bucket[*job.Job]) error { return job.NewStore(buc).Put(j) })
		}).AsTask(),
		0,
	)

	rep, err := reporter.OpenJSON(reportPath)
	if err != nil {
		return nil, err
	}

	sin := task.Connect(
		js.AsTask(),
		sink.ElementWise(func(ctx context.Context, j *job.Job) error { return rep.Write(j) }).AsTask(),
		0,
	)
	sin.Defer(func() { rep.Close() })

	return flow.FromTask(sin), nil
}

func start(ctx context.Context, paths []string, reportPath string) error {
	cacheDB, err := database.Open(bucket)
	if err != nil {
		return err
	}
	defer cacheDB.Close()

	fl, err := downloadFiles(paths, reportPath, cacheDB)
	if err != nil {
		return err
	}
	return fl.Run(ctx)
}
