package drawin

type RetryConfig struct {
	ShouldRetry func(*Report) bool
}

var DefaultRetryConfig = &RetryConfig{
	ShouldRetry: func(rep *Report) bool {
		if rep.Result != FailedResult {
			return false
		}
		return !rep.Failure.Permanent
	},
}
