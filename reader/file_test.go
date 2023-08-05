package reader_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/carbo/pipe"
	"github.com/hiroara/carbo/taskfn"
	"github.com/hiroara/drawin/reader"
)

func TestRead(t *testing.T) {
	t.Parallel()

	readFiles := taskfn.SliceToSlice(pipe.FromFn(reader.Read).AsTask())
	out, err := readFiles(context.Background(), []string{"file.go"})
	require.NoError(t, err)
	assert.Contains(t, out, "package reader")
}
