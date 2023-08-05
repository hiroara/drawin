package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/drawin/database"
)

func TestBucket(t *testing.T) {
	t.Parallel()

	bucket := []byte("test-bucket")
	key1 := []byte("item1")

	db, err := database.Open(bucket)
	require.NoError(t, err)
	err = database.Update(db, func(buc *database.Bucket[*entry]) error {
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
