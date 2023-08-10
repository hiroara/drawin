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
	outStr      string
	reportPath  string
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
			return download(downloadFS.Args(), outStr, reportPath, concurrency)
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
	downloadFS.StringVar(&outStr, "out", "drawin-out", `Output configuration with the format "<type>=<path>".
Available output types: directory|store

Also, it is possible to specify only "<path>" (without "<type">).
In this case, it is interpreted as a shorthand of "directory=<path>".`)
	downloadFS.StringVar(&reportPath, "report", "-", "Path to the file that a download report is written (\"-\" means STDOUT).")
	downloadFS.IntVar(&concurrency, "concurrency", 6, "Number of concurrent connectinos.")
}
