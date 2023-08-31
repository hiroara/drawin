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
)

type dummyClient struct {
	mock.Mock
}

func (cli *dummyClient) Download(ctx context.Context, j *job.Job) (*drawin.Report, error) {
	args := cli.Mock.Called(ctx, j)
	return args.Get(0).(*drawin.Report), args.Error(1)
}

var urls = []string{
	"https://example.com/assets/image1.jpg",
	"https://example.com/assets/image2.jpg",
}

func setupClient(cli *dummyClient, urls []string) {
	jobExpectation1 := mock.MatchedBy(func(j *job.Job) bool {
		return j.URL == urls[0]
	})

	jobExpectation2 := mock.MatchedBy(func(j *job.Job) bool {
		return j.URL == urls[1]
	})

	cli.On("Download", mock.Anything, jobExpectation1).Return(drawin.Downloaded(&job.Job{Name: "image1.jpg", URL: urls[0]}, 1024), nil).Once()
	cli.On("Download", mock.Anything, jobExpectation2).Return(drawin.Downloaded(&job.Job{Name: "image2.jpg", URL: urls[1]}, 1024), nil).Once()
}

func TestDownloader(t *testing.T) {
	t.Parallel()

	cli := new(dummyClient)
	d, err := drawin.New(cli)
	require.NoError(t, err)

	urlsC := make(chan string, 2)
	urlsC <- urls[0]
	urlsC <- urls[1]
	close(urlsC)

	setupClient(cli, urls)

	ctx := context.Background()

	out := make(chan *drawin.Report, 2)

	require.NoError(t, d.Run(ctx, urlsC, out))

	results := make([]string, 0)
	for rep := range out {
		results = append(results, rep.URL)
	}
	assert.ElementsMatch(t, urls, results)

	cli.AssertExpectations(t)
}

func TestDownloaderAsTask(t *testing.T) {
	t.Parallel()

	cli := new(dummyClient)
	d, err := drawin.New(cli)
	require.NoError(t, err)

	dfn := taskfn.SliceToSlice(d.AsTask())

	setupClient(cli, urls)

	reps, err := dfn(context.Background(), urls)
	require.NoError(t, err)

	assert.Len(t, reps, len(urls))
}
