package database_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/drawin/database"
)

func TestBucket(t *testing.T) {
	t.Parallel()

	key1 := []byte("item1")

	path := filepath.Join(t.TempDir(), "test.db")
	sdb, err := openSingleDB(path)
	require.NoError(t, err)
	err = sdb.Update(func(buc *database.Bucket[*entry]) error {
		e, err := buc.Get(key1)
		require.NoError(t, err)
		assert.Nil(t, e)

		err = buc.Put(key1, &entry{Name: "entry1"})
		require.NoError(t, err)

		e, err = buc.Get(key1)
		require.NoError(t, err)
		if assert.NotNil(t, e) {
			assert.Equal(t, "entry1", e.Name)
		}

		return nil
	})
	require.NoError(t, err)
}
