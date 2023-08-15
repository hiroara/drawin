package client

import "github.com/hiroara/drawin/reporter"

type RetryConfig struct {
	ShouldRetry func(*reporter.Report) bool
}

var DefaultRetryConfig = &RetryConfig{
	ShouldRetry: func(rep *reporter.Report) bool {
		if rep.Result != reporter.Failed {
			return false
		}
		return !rep.Failure.Permanent
	},
}

func (cfg *RetryConfig) shouldRetry(rep *reporter.Report) bool {
	if cfg == nil || cfg.ShouldRetry == nil {
		return DefaultRetryConfig.ShouldRetry(rep)
	}
	return cfg.ShouldRetry(rep)
}
