package drawin

import "github.com/hiroara/drawin/job"

type Report struct {
	job.Job       `json:"job"`
	Result        Result   `json:"result"`
	ContentLength int64    `json:"contentLength,omitempty"`
	Failure       *Failure `json:"failure,omitempty"`
}

type Failure struct {
	Permanent bool   `json:"permanent"`
	Error     string `json:"error"`
}

func DownloadedReport(j *job.Job, contentLength int64) *Report {
	return &Report{Job: *j, Result: DownloadedResult, ContentLength: contentLength}
}

func CachedReport(j *job.Job, contentLength int64) *Report {
	return &Report{Job: *j, Result: CachedResult, ContentLength: contentLength}
}

func FailedReport(j *job.Job, err error, permanent bool) *Report {
	return &Report{Job: *j, Result: FailedResult, Failure: &Failure{Permanent: permanent, Error: err.Error()}}
}

type Result string

var (
	DownloadedResult Result = "downloaded"
	CachedResult     Result = "cached"
	FailedResult     Result = "failed"
)
