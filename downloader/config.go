package downloader

type config struct {
	concurrency int
	batchSize      int
}

func newConfig(opts ...Option) *config {
	cfg := &config{
		concurrency: 4,
		batchSize:   128,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}
