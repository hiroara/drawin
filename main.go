package main

import (
	"context"
	"flag"
	"log"
)

var bucket = []byte("drawin")

var (
	outdir      string
	concurrency int
)

func main() {
	flag.Parse()

	if err := start(context.Background(), flag.Args()); err != nil {
		log.Fatal(err)
	}
}

func init() {
	flag.StringVar(&outdir, "outdir", "drawin-out", "Path to the directory to download files")
	flag.IntVar(&concurrency, "concurrency", 6, "Number of concurrent connectinos")
}
