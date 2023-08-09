package reader_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/source"
	"github.com/hiroara/carbo/taskfn"
	"github.com/hiroara/drawin/reader"
)

func TestRead(t *testing.T) {
	t.Parallel()

	readFiles := taskfn.SourceToSlice(source.FromFn(func(ctx context.Context, out chan<- string) error {
		return reader.Read(ctx, "file.go", out)
	}))
	out, err := readFiles(context.Background())
	require.NoError(t, err)
	assert.Contains(t, out, "package reader")
}
