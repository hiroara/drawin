package drawin

type Option func(*config)

func WithConcurrency(concurrency int) Option {
	return func(cfg *config) {
		cfg.concurrency = concurrency
	}
}

type config struct {
	concurrency int
	batchSize   int
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
