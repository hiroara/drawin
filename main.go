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
var readFS *flag.FlagSet

func main() {
	flag.Parse()

	com := flag.Arg(0)

	if com == "" {
	}

	switch com {
	case "download":
		downloadFS.Parse(flag.Args()[1:])
	case "read":
		readFS.Parse(flag.Args()[1:])
		if readFS.NArg() < 1 {
			readFS.Usage()
			os.Exit(1)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}

	reg := registry.New()

	reg.Register(
		"download",
		flow.NewFactory(func() (*flow.Flow, error) {
			return download(downloadFS.Args(), outStr, reportPath, concurrency)
		}),
	)

	reg.Register(
		"read",
		flow.NewFactory(func() (*flow.Flow, error) {
			return read(readFS.Arg(0), readFS.Args()[1:])
		}),
	)

	if err := reg.Run(context.Background(), com); err != nil {
		log.Fatal(err)
	}
}

func init() {
	flag.Usage = usage(flag.CommandLine, os.Args[0], "<command>", "Available commands: download|read")

	downloadFS = flag.NewFlagSet("download", flag.ExitOnError)
	downloadFS.Usage = usage(downloadFS, os.Args[0], "download [<list of URLs>...]", `Download files from URLs listed in the passed files.
Passing "-" as a positional argument means reading list of URLs from STDIN.`)
	downloadFS.StringVar(&outStr, "out", "drawin-out", `Output configuration with the format "<type>=<path>".
Available output types: directory|store

Also, it is possible to specify only "<path>" (without "<type">).
In this case, it is interpreted as a shorthand of "directory=<path>".`)
	downloadFS.StringVar(&reportPath, "report", "-", "Path to the file that a download report is written (\"-\" means STDOUT).")
	downloadFS.IntVar(&concurrency, "concurrency", 6, "Number of concurrent connectinos.")

	readFS = flag.NewFlagSet("read", flag.ExitOnError)
	readFS.Usage = usage(readFS, os.Args[0], "read <path to store> [<URL>...]", "_TBD_")
}

func usage(fs *flag.FlagSet, prog, command, desc string) func() {
	return func() {
		fmt.Fprintf(fs.Output(), "Usage:\n  %s [<global option>...] %s [<command option>...]\n", prog, command)
		if desc != "" {
			fmt.Fprintf(fs.Output(), "\n%s\n\n", desc)
		}
		fs.PrintDefaults()
	}
}
