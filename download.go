package main

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/hiroara/drawin/client"
	"github.com/hiroara/drawin/downloader"
	"github.com/hiroara/drawin/reader"
	"github.com/hiroara/drawin/reporter"
)

func start(ctx context.Context, paths []string, reportPath string) error {
	cli := client.New(outdir)
	if err := cli.CreateDir(); err != nil {
		return err
	}

	d, err := downloader.New(cli)
	if err != nil {
		return err
	}
	grp, ctx := errgroup.WithContext(ctx)

	urls := make(chan string)

	rep, err := reporter.OpenJSON(reportPath)
	if err != nil {
		return err
	}

	grp.Go(func() error {
		defer close(urls)
		return reader.Read(ctx, paths, urls)
	})

	grp.Go(func() error {
		defer d.Close()
		return d.Run(ctx, urls)
	})

	grp.Go(func() error {
		for j := range d.Output() {
			if err := rep.Write(j); err != nil {
				return err
			}
		}
		return nil
	})

	return grp.Wait()
}
