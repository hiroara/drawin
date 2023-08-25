package downloader

import (
	"context"
	"os"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/sink"
	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/task"
	"github.com/hiroara/drawin/database"
	"github.com/hiroara/drawin/downloader/report"
	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/marshal"
)

type Downloader struct {
	client    Client
	cache     *database.SingleDB[*job.Job]
	cacheFile *os.File
	config    *config
}

type Report = report.Report

type Client interface {
	Download(ctx context.Context, j *job.Job) (*Report, error)
}

var cacheBucket = []byte("drawin-cache")

func New(cli Client, opts ...Option) (*Downloader, error) {
	f, err := os.CreateTemp("", "drawin-*.db")
	if err != nil {
		return nil, err
	}

	cacheDB, err := database.Open(f.Name(), nil)
	if err != nil {
		return nil, err
	}
	cacheSDB, err := database.Single[*job.Job](cacheDB, cacheBucket, marshal.Msgpack[*job.Job]())

	return &Downloader{
		client:    cli,
		cache:     cacheSDB,
		cacheFile: f,
		config:    newConfig(opts...),
	}, nil
}

func (d *Downloader) Run(ctx context.Context, urls <-chan string, out chan<- *Report) error {
	f, err := d.downloadFlow(urls, out)
	if err != nil {
		return err
	}

	return f.Run(ctx)
}

func (d *Downloader) downloadFlow(urls <-chan string, out chan<- *Report) (*flow.Flow, error) {
	src := source.FromChan(urls)
	urlBatches := task.Connect(
		src.AsTask(),
		pipe.Batch[string](d.config.batchSize).AsTask(),
		0,
	)
	jobBatches := task.Connect(
		urlBatches.AsTask(),
		pipe.Map(func(ctx context.Context, urls []string) ([]*job.Job, error) {
			jobs := make([]*job.Job, 0, len(urls))
			err := d.cache.Update(func(buc *database.Bucket[*job.Job]) error {
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
	jobs := task.Connect(
		jobBatches.AsTask(),
		pipe.FlattenSlice[*job.Job]().AsTask(),
		0,
	)
	reps := task.Connect(
		jobs,
		pipe.Map(d.client.Download).Concurrent(d.config.concurrency).AsTask(),
		0,
	)

	sin := task.Connect(
		reps.AsTask(),
		sink.ToChan(out).AsTask(),
		0,
	)

	return flow.FromTask(sin), nil
}

func (d *Downloader) Close() error {
	d.cacheFile.Close()
	os.Remove(d.cacheFile.Name())
	return nil
}

func (d *Downloader) AsTask() task.Task[string, *Report] {
	return task.FromFn(d.Run)
}
