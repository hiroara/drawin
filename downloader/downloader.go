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
	"github.com/hiroara/drawin/job"
)

type Downloader struct {
	client      Client
	cache       *database.SingleDB[*job.Job]
	cacheFile   *os.File
	concurrency int
	out         chan *job.Job
}

type config struct {
	concurrency int
	buffer      int
}

type Client interface {
	Download(ctx context.Context, j *job.Job) error
}

var cacheBucket = []byte("drawin-cache")

func New(cli Client) (*Downloader, error) {
	cfg := &config{
		concurrency: 4,
		buffer:      64,
	}

	f, err := os.CreateTemp("", "drawin-*.db")
	if err != nil {
		return nil, err
	}

	cacheDB, err := database.Open(f.Name())
	if err != nil {
		return nil, err
	}
	cacheSDB, err := database.Single[*job.Job](cacheDB, cacheBucket)

	return &Downloader{
		client:      cli,
		cache:       cacheSDB,
		cacheFile:   f,
		concurrency: cfg.concurrency,
		out:         make(chan *job.Job, cfg.buffer),
	}, nil
}

func (d *Downloader) Output() <-chan *job.Job {
	return d.out
}

func (d *Downloader) Run(ctx context.Context, urls <-chan string) error {
	f, err := d.downloadFlow(urls)
	if err != nil {
		return err
	}

	return f.Run(ctx)
}

func (d *Downloader) downloadFlow(urls <-chan string) (*flow.Flow, error) {
	src := source.FromFn(func(ctx context.Context, out chan<- string) error {
		for url := range urls {
			if err := task.Emit(ctx, out, url); err != nil {
				return err
			}
		}
		return nil
	})
	urlBatches := task.Connect(
		src.AsTask(),
		pipe.Batch[string](32).AsTask(),
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
	jobs = task.Connect(
		jobs.AsTask(),
		pipe.Tap(d.client.Download).Concurrent(d.concurrency).AsTask(),
		0,
	)
	jobs = task.Connect(
		jobs.AsTask(),
		pipe.Tap(func(ctx context.Context, j *job.Job) error {
			return d.cache.Update(func(buc *database.Bucket[*job.Job]) error { return job.NewStore(buc).Put(j) })
		}).AsTask(),
		0,
	)

	sin := task.Connect(
		jobs.AsTask(),
		sink.ElementWise(func(ctx context.Context, j *job.Job) error {
			return task.Emit(ctx, d.out, j)
		}).AsTask(),
		0,
	)

	return flow.FromTask(sin), nil
}

func (d *Downloader) Close() error {
	close(d.out)
	d.cacheFile.Close()
	os.Remove(d.cacheFile.Name())
	return nil
}
