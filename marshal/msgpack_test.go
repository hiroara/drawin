package marshal_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/drawin/marshal"
)

type entry struct {
	Name string
}

func TestMsgpack(t *testing.T) {
	t.Parallel()

	sp := marshal.Msgpack[*entry]()

	ent1 := &entry{Name: "entry1"}

	bs, err := sp.Marshal(ent1)
	require.NoError(t, err)

	// fixmap with 1 length
	assert.Equal(t, "81", fmt.Sprintf("%x", bs[:1]))

	v, err := sp.Unmarshal(bs)
	require.NoError(t, err)

	assert.Equal(t, ent1, v)
}
