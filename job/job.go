package job

type Job struct {
	URL        string
	Name       string
	Body       []byte
	Downloaded bool
}
