package reporter

import "github.com/hiroara/drawin/job"

type Report struct {
	job.Job       `json:"job"`
	Result        Result `json:"result"`
	ContentLength int64  `json:"contentLength,omitempty"`
	Error         string `json:"error,omitempty"`
}

func DownloadedReport(j *job.Job, contentLength int64) *Report {
	return &Report{Job: *j, Result: Downloaded, ContentLength: contentLength}
}

func CachedReport(j *job.Job, contentLength int64) *Report {
	return &Report{Job: *j, Result: Cached, ContentLength: contentLength}
}

func FailedReport(j *job.Job, err error) *Report {
	return &Report{Job: *j, Result: Failed, Error: err.Error()}
}

type Result string

var (
	Downloaded Result = "downloaded"
	Cached     Result = "cached"
	Failed     Result = "failed"
)
