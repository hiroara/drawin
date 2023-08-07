package main

import (
	"context"
	"flag"
	"log"
)

var bucket = []byte("drawin")

var (
	outdir      string
	reportPath  string
	concurrency int
)

func main() {
	flag.Parse()

	if err := start(context.Background(), flag.Args(), reportPath); err != nil {
		log.Fatal(err)
	}
}

func init() {
	flag.StringVar(&outdir, "outdir", "drawin-out", "path to the directory to download files")
	flag.StringVar(&reportPath, "report", "-", "path to the file that a download report is written (defaults to `-` that means STDOUT)")
	flag.IntVar(&concurrency, "concurrency", 6, "number of concurrent connectinos")
}
