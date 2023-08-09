package downloader_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/drawin/downloader"
	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/reporter"
)

type dummyClient struct {
	mock.Mock
}

func (cli *dummyClient) Download(ctx context.Context, j *job.Job) (*reporter.Report, error) {
	args := cli.Mock.Called(ctx, j)
	return args.Get(0).(*reporter.Report), args.Error(1)
}

func TestDownloader(t *testing.T) {
	t.Parallel()

	urls := []string{
		"https://example.com/assets/image1.jpg",
		"https://example.com/assets/image2.jpg",
	}

	cli := new(dummyClient)
	out := make(chan *reporter.Report, 2)

	d, err := downloader.New(cli, out)
	require.NoError(t, err)

	urlsC := make(chan string, 2)
	urlsC <- urls[0]
	urlsC <- urls[1]
	close(urlsC)

	jobExpectation1 := mock.MatchedBy(func(j *job.Job) bool {
		return j.URL == urls[0]
	})

	jobExpectation2 := mock.MatchedBy(func(j *job.Job) bool {
		return j.URL == urls[1]
	})

	cli.On("Download", mock.Anything, jobExpectation1).Return(reporter.DownloadedReport(&job.Job{Name: "image1.jpg", URL: urls[0]}, 1024), nil).Once()
	cli.On("Download", mock.Anything, jobExpectation2).Return(reporter.DownloadedReport(&job.Job{Name: "image2.jpg", URL: urls[1]}, 1024), nil).Once()

	ctx := context.Background()

	require.NoError(t, d.Run(ctx, urlsC))

	results := make([]string, 0)
	for rep := range out {
		results = append(results, rep.URL)
	}
	assert.ElementsMatch(t, urls, results)

	cli.AssertExpectations(t)
}
