package drawin

import (
	"github.com/hiroara/drawin/handler"
)

type Option func(*config)

func WithConcurrency(concurrency int) Option {
	return func(cfg *config) {
		cfg.concurrency = concurrency
	}
}

func WithHandlers(hs ...handler.Handler) Option {
	return func(cfg *config) {
		cfg.handlers = hs
	}
}

type config struct {
	concurrency int
	batchSize   int
	handlers    []handler.Handler
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
