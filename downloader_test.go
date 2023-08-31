package drawin_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/taskfn"
	"github.com/hiroara/drawin"
	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/output"
)

type dummyHandler struct {
	mock.Mock
}

func (h *dummyHandler) Match(j *job.Job) bool {
	args := h.Mock.Called(j)
	return args.Bool(0)
}

func (h *dummyHandler) ShouldRetry(err error) bool {
	args := h.Mock.Called(err)
	return args.Bool(0)
}

func (h *dummyHandler) Get(ctx context.Context, j *job.Job) ([]byte, error) {
	args := h.Mock.Called(ctx, j)
	return args.Get(0).([]byte), args.Error(1)
}

var urls = []string{
	"https://example.com/assets/image1.jpg",
	"https://example.com/assets/image2.jpg",
}

func setupHandler(urls []string) *dummyHandler {
	h := &dummyHandler{}

	jobExpectation1 := mock.MatchedBy(func(j *job.Job) bool {
		return j.URL == urls[0]
	})

	jobExpectation2 := mock.MatchedBy(func(j *job.Job) bool {
		return j.URL == urls[1]
	})

	h.On("Match", jobExpectation1).Return(true)
	h.On("Match", jobExpectation2).Return(true)

	h.On("Get", mock.Anything, jobExpectation1).Return([]byte("result1"), nil)
	h.On("Get", mock.Anything, jobExpectation2).Return([]byte("result2"), nil)

	return h
}

func TestDownloader(t *testing.T) {
	t.Parallel()

	outDir := output.NewDirectory(t.TempDir())
	h := setupHandler(urls)

	d, err := drawin.NewDownloader(outDir, drawin.WithHandlers(h))
	require.NoError(t, err)

	urlsC := make(chan string, 2)
	urlsC <- urls[0]
	urlsC <- urls[1]
	close(urlsC)

	ctx := context.Background()

	out := make(chan *drawin.Report, 2)

	require.NoError(t, d.Run(ctx, urlsC, out))

	results := make([]string, 0)
	for rep := range out {
		results = append(results, rep.URL)
	}
	assert.ElementsMatch(t, urls, results)

	h.AssertExpectations(t)
}

func TestDownloaderAsTask(t *testing.T) {
	t.Parallel()

	outDir := output.NewDirectory(t.TempDir())
	h := setupHandler(urls)

	d, err := drawin.NewDownloader(outDir, drawin.WithHandlers(h))
	require.NoError(t, err)

	dfn := taskfn.SliceToSlice(d.AsTask())

	reps, err := dfn(context.Background(), urls)
	require.NoError(t, err)

	assert.Len(t, reps, len(urls))
}
