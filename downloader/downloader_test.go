package downloader_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

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
	d, err := downloader.New(cli)
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
	grp, ctx := errgroup.WithContext(ctx)

	grp.Go(func() error {
		defer d.Close()
		return d.Run(ctx, urlsC)
	})

	grp.Go(func() error {
		results := make([]string, 0)
		for rep := range d.Output() {
			results = append(results, rep.URL)
		}
		assert.ElementsMatch(t, urls, results)
		return nil
	})

	require.NoError(t, grp.Wait())

	cli.AssertExpectations(t)
}
