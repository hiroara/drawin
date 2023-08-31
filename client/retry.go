package client

import "github.com/hiroara/drawin"

type RetryConfig struct {
	ShouldRetry func(*drawin.Report) bool
}

var DefaultRetryConfig = &RetryConfig{
	ShouldRetry: func(rep *drawin.Report) bool {
		if rep.Result != drawin.FailedResult {
			return false
		}
		return !rep.Failure.Permanent
	},
}

func (cfg *RetryConfig) shouldRetry(rep *drawin.Report) bool {
	if cfg == nil || cfg.ShouldRetry == nil {
		return DefaultRetryConfig.ShouldRetry(rep)
	}
	return cfg.ShouldRetry(rep)
}
