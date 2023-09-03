package database_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"

	"github.com/hiroara/drawin/internal/database"
	"github.com/hiroara/drawin/marshal"
)

type entry struct {
	Name string
}

var bucket = []byte("test-bucket")

func openSingleDB(path string) (*database.SingleDB[*entry], error) {
	return database.OpenSingle[*entry](path, bucket, marshal.Msgpack[*entry](), true)
}

func TestSingle(t *testing.T) {
	t.Parallel()

	db, err := openDB(filepath.Join(t.TempDir(), "test.db"), true)
	require.NoError(t, err)
	sdb, err := database.Single[*entry](db, bucket, marshal.Msgpack[*entry]())
	require.NoError(t, err)
	assert.NotNil(t, sdb)
	require.NoError(t, sdb.Close())
}

func TestSingleDBUpdateAndView(t *testing.T) {
	t.Parallel()

	key := []byte("test-key")
	ent := &entry{Name: "entry1"}

	path := filepath.Join(t.TempDir(), "test.db")
	sdb, err := openSingleDB(path)
	require.NoError(t, err)
	defer sdb.Close()

	err = sdb.Update(func(buc *database.Bucket[*entry]) error {
		assert.NotNil(t, buc)
		return buc.Put(key, ent)
	})
	require.NoError(t, err)

	err = sdb.View(func(buc *database.Bucket[*entry]) error {
		assert.NotNil(t, buc)
		v, err := buc.Get(key)
		require.NoError(t, err)
		assert.Equal(t, ent, v)
		return nil
	})
	require.NoError(t, err)

	err = sdb.View(func(buc *database.Bucket[*entry]) error {
		assert.NotNil(t, buc)
		v, err := buc.Get(key)
		require.NoError(t, err)
		assert.Equal(t, ent, v)
		return nil
	})
	require.NoError(t, err)

	require.NoError(t, sdb.Close())

	db, err := openDB(path, true)
	require.NoError(t, err)

	err = db.View(func(tx *bolt.Tx) error {
		bs := tx.Bucket(bucket).Get(key)
		v, err := marshal.Msgpack[*entry]().Unmarshal(bs)
		require.NoError(t, err)
		assert.Equal(t, ent, v)
		return nil
	})
	require.NoError(t, err)
}
