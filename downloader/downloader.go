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
	"github.com/hiroara/drawin/reporter"
	"golang.org/x/sync/errgroup"
)

type Downloader struct {
	client      Client
	cache       *database.SingleDB[*job.Job]
	cacheFile   *os.File
	concurrency int
}

type config struct {
	concurrency int
	buffer      int
}

type Client interface {
	Download(ctx context.Context, j *job.Job) (*reporter.Report, error)
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
	}, nil
}

func (d *Downloader) Run(ctx context.Context, urls <-chan string, out chan<- *reporter.Report) error {
	f, err := d.downloadFlow(urls, out)
	if err != nil {
		return err
	}

	return f.Run(ctx)
}

func (d *Downloader) downloadFlow(urls <-chan string, out chan<- *reporter.Report) (*flow.Flow, error) {
	src := source.FromChan(urls)
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
	reps := task.Connect(
		jobs,
		pipe.Map(d.client.Download).Concurrent(d.concurrency).AsTask(),
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

func (d *Downloader) AsPipe() pipe.Pipe[string, *reporter.Report] {
	return pipe.FromFn(func(ctx context.Context, urls <-chan string, out chan<- *reporter.Report) error {
		m := make(chan *reporter.Report)
		grp, ctx := errgroup.WithContext(ctx)
		grp.Go(func() error { return d.Run(ctx, urls, m) })
		grp.Go(func() error {
			for rep := range m {
				if err := task.Emit(ctx, out, rep); err != nil {
					return err
				}
			}
			return nil
		})
		return grp.Wait()
	})
}

func (d *Downloader) AsTask() task.Task[string, *reporter.Report] {
	return d.AsPipe()
}
