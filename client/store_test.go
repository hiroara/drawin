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
	"github.com/hiroara/drawin/reporter"
)

func TestStoreOutput(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "out", "test.db")
	db, err := database.Open(path)
	require.NoError(t, err)
	defer db.Close()

	out := client.NewStore(db)
	require.NoError(t, out.Prepare())

	j := &job.Job{Name: "file1.txt", URL: "https://example.com/dir/file1.txt"}
	data := []byte("test value")

	rep, err := out.Get(j)
	require.NoError(t, err)
	assert.Nil(t, rep)

	require.NoError(t, out.Add(reporter.DownloadedReport(j, int64(len(data))), data))

	require.NoError(t, db.View(func(tx *bolt.Tx) error {
		imgs := tx.Bucket([]byte("images"))
		v := imgs.Get([]byte(j.URL))
		assert.Equal(t, []byte("test value"), v)

		reps := tx.Bucket([]byte("reports"))
		bs := reps.Get([]byte(j.URL))
		rep, err := marshal.Gob[*reporter.Report]().Unmarshal(bs)
		require.NoError(t, err)
		assert.Equal(t, "file1.txt", rep.Name)
		return nil
	}))

	rep, err = out.Get(j)
	require.NoError(t, err)
	if assert.NotNil(t, rep) {
		assert.Equal(t, reporter.Cached, rep.Result)
		assert.Equal(t, int64(len(data)), rep.ContentLength)
	}
}
