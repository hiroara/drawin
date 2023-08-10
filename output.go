package main

import (
	"errors"
	"fmt"
	"strings"
)

type output struct {
	typ  outputType
	path string
}

type outputType string

var (
	directoryType outputType = "directory"
	storeType     outputType = "store"
)

var availableOutputTypes = []outputType{directoryType, storeType}

var errNoMatchingOutputType = errors.New("unknown output type is specified")

func parseOutput(s string) (*output, error) {
	ss := strings.SplitN(s, "=", 2)
	if len(ss) == 0 {
		return nil, fmt.Errorf("%w: <empty>", errNoMatchingOutputType, s)
	}

	if len(ss) == 1 {
		return &output{typ: directoryType, path: ss[0]}, nil
	}
	t := ss[0]
	v := ss[1]

	match := false
	for _, c := range availableOutputTypes {
		match = match || t == string(c)
	}

	if !match {
		return nil, fmt.Errorf("%w: %s", errNoMatchingOutputType, s)
	}

	return &output{typ: outputType(t), path: v}, nil
}
