package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hiroara/carbo/flow"
	"github.com/hiroara/carbo/registry"
)

var bucket = []byte("drawin")

var (
	outdir      string
	reportPath  string
	useStore    bool
	concurrency int
)

var downloadFS *flag.FlagSet

func main() {
	flag.Parse()

	com := flag.Arg(0)

	if com == "download" {
		downloadFS.Parse(flag.Args()[1:])
	}

	reg := registry.New()

	reg.Register(
		"download",
		flow.NewFactory(func() (*flow.Flow, error) {
			return download(downloadFS.Args(), outdir, reportPath, useStore, concurrency)
		}),
	)

	if err := reg.Run(context.Background(), com); err != nil {
		log.Fatal(err)
	}
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [<global option>, ...] <subcommand> [<command option>, ...]\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Available subcommands: download\n")
		flag.PrintDefaults()
	}

	downloadFS = flag.NewFlagSet(fmt.Sprintf("%s download", os.Args[0]), flag.ExitOnError)
	downloadFS.StringVar(&outdir, "outdir", "drawin-out", "path to the directory to download files")
	downloadFS.BoolVar(&useStore, "store", false, "enable store download mode")
	downloadFS.StringVar(&reportPath, "report", "-", "path to the file that a download report is written (\"-\" means STDOUT)")
	downloadFS.IntVar(&concurrency, "concurrency", 6, "number of concurrent connectinos")
}
