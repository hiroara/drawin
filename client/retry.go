package client

import (
	"github.com/hiroara/drawin/downloader/report"
)

type RetryConfig struct {
	ShouldRetry func(*report.Report) bool
}

var DefaultRetryConfig = &RetryConfig{
	ShouldRetry: func(rep *report.Report) bool {
		if rep.Result != report.FailedResult {
			return false
		}
		return !rep.Failure.Permanent
	},
}

func (cfg *RetryConfig) shouldRetry(rep *report.Report) bool {
	if cfg == nil || cfg.ShouldRetry == nil {
		return DefaultRetryConfig.ShouldRetry(rep)
	}
	return cfg.ShouldRetry(rep)
}
