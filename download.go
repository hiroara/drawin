package main

import (
	"context"
	"path/filepath"

	"golang.org/x/sync/errgroup"

	"github.com/hiroara/drawin/client"
	"github.com/hiroara/drawin/database"
	"github.com/hiroara/drawin/downloader"
	"github.com/hiroara/drawin/reader"
	"github.com/hiroara/drawin/reporter"
)

func start(ctx context.Context, paths []string, outdir, reportPath string, useStore bool, concurrency int) error {
	var out client.Output
	if useStore {
		db, err := database.Open(filepath.Join(outdir, "drawin.db"))
		if err != nil {
			return err
		}
		defer db.Close()

		out = client.NewStore(db)
	} else {
		out = client.NewDirectory(outdir)
	}

	cli, err := client.Build(out)
	if err != nil {
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
