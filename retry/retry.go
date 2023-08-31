package retry

import "github.com/hiroara/drawin/report"

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
