package job

type Job struct {
	URL           string `json:"url"`
	Name          string `json:"name"`
	Body          []byte `json:"-"`
	ContentLength int64  `json:"contentLength"`
	Action        Action `json:"action"`
}

type Action string

var (
	DownloadAction Action = "download"
	CacheAction    Action = "cache"
)
