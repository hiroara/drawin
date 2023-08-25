package client

import (
	"github.com/hiroara/drawin/downloader"
	"github.com/hiroara/drawin/downloader/report"
)

type RetryConfig struct {
	ShouldRetry func(*downloader.Report) bool
}

var DefaultRetryConfig = &RetryConfig{
	ShouldRetry: func(rep *downloader.Report) bool {
		if rep.Result != report.FailedResult {
			return false
		}
		return !rep.Failure.Permanent
	},
}

func (cfg *RetryConfig) shouldRetry(rep *downloader.Report) bool {
	if cfg == nil || cfg.ShouldRetry == nil {
		return DefaultRetryConfig.ShouldRetry(rep)
	}
	return cfg.ShouldRetry(rep)
}
