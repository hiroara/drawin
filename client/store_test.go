package client_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"

	"github.com/hiroara/carbo/marshal"
	"github.com/hiroara/drawin/client"
	"github.com/hiroara/drawin/database"
	"github.com/hiroara/drawin/job"
)

func TestStoreOutput(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "out", "test.db")
	db, err := database.Open(path)
	require.NoError(t, err)
	defer db.Close()

	out := client.NewStore(db)
	require.NoError(t, out.Prepare())

	j := &job.Job{Name: "file1.txt"}

	ok, err := out.Check(j)
	require.NoError(t, err)
	assert.False(t, ok)

	require.NoError(t, out.Add(j, []byte("test value")))

	require.NoError(t, db.View(func(tx *bolt.Tx) error {
		imgs := tx.Bucket([]byte("images"))
		v := imgs.Get([]byte(j.Name))
		assert.Equal(t, []byte("test value"), v)

		reps := tx.Bucket([]byte("reports"))
		bs := reps.Get([]byte(j.Name))
		j, err := marshal.Gob[*job.Job]().Unmarshal(bs)
		require.NoError(t, err)
		assert.Equal(t, "file1.txt", j.Name)
		return nil
	}))

	ok, err = out.Check(j)
	require.NoError(t, err)
	assert.True(t, ok)
}
