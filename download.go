package main

import (
	"context"
	"log"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"

	"github.com/hiroara/drawin/database"
	"github.com/hiroara/drawin/downloader"
	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/reader"
)

func debug(ctx context.Context, j *job.Job) error {
	action := "downloaded"
	if !j.Downloaded {
		action = "cache"
	}
	log.Printf("%s => %s (%s)", j.URL, j.Name, action)
	return nil
}

func downloadFiles(paths []string, db *database.DB) (*flow.Flow, error) {
	d := downloader.New(outdir)
	if err := d.CreateDir(); err != nil {
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
			err := database.Update(db, func(buc *database.Bucket[*job.Job]) error {
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
		0,
	)
	js := task.Connect(
		jbs.AsTask(),
		pipe.FlattenSlice[*job.Job]().AsTask(),
		0,
	)
	bs := task.Connect(
		js.AsTask(),
		pipe.Tap(d.Download).Concurrent(concurrency).AsTask(),
		0,
	)
	sin := task.Connect(
		bs.AsTask(),
		sink.ElementWise(debug).AsTask(),
		0,
	)

	return flow.FromTask(sin), nil
}

func start(ctx context.Context, paths []string) error {
	db, err := database.Open(bucket)
	if err != nil {
		return err
	}
	defer db.Close()

	fl, err := downloadFiles(paths, db)
	if err != nil {
		return err
	}
	return fl.Run(ctx)
}
