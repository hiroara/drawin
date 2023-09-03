package store_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"

	"github.com/hiroara/drawin/internal/database"
	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/marshal"
	"github.com/hiroara/drawin/report"
	"github.com/hiroara/drawin/store"
)

func TestStore(t *testing.T) {
	t.Parallel()

	buildStore := func(t *testing.T) (*database.DB, *store.Store) {
		path := filepath.Join(t.TempDir(), "out", "test.db")
		db, err := database.Open(path, true)
		require.NoError(t, err)

		s := store.New(db)
		require.NoError(t, s.Initialize())
		return db, s
	}

	j := &job.Job{Name: "file1.txt", URL: "https://example.com/dir/file1.txt"}

	t.Run("NoDataCase", func(t *testing.T) {
		t.Parallel()

		_, s := buildStore(t)

		rep, err := s.Get(j)
		require.NoError(t, err)
		assert.Nil(t, rep)
	})

	t.Run("DownloadedReportCase", func(t *testing.T) {
		t.Parallel()

		data := []byte("test value")

		db, s := buildStore(t)

		require.NoError(t, s.Add(report.Downloaded(j, int64(len(data))), data))

		require.NoError(t, db.View(func(tx *bolt.Tx) error {
			imgs := tx.Bucket([]byte("images"))
			v := imgs.Get([]byte(j.URL))
			assert.Equal(t, data, v)

			reps := tx.Bucket([]byte("reports"))
			bs := reps.Get([]byte(j.URL))
			rep, err := marshal.Msgpack[*report.Report]().Unmarshal(bs)
			require.NoError(t, err)
			assert.Equal(t, "file1.txt", rep.Name)
			return nil
		}))

		rep, err := s.Get(j)
		require.NoError(t, err)
		if assert.NotNil(t, rep) {
			assert.Equal(t, report.DownloadedResult, rep.Result)
			assert.Equal(t, int64(len(data)), rep.ContentLength)
		}

		bs, err := s.Read(rep)
		require.NoError(t, err)
		if assert.NotNil(t, bs) {
			assert.Equal(t, data, bs)
		}
	})

	t.Run("FailedReportCase", func(t *testing.T) {
		t.Parallel()

		_, s := buildStore(t)

		require.NoError(t, s.Add(report.Failed(j, errors.New("test error"), true), nil))

		rep, err := s.Get(j)
		require.NoError(t, err)
		if assert.NotNil(t, rep) {
			assert.Equal(t, report.FailedResult, rep.Result)
		}

		bs, err := s.Read(rep)
		require.NoError(t, err)
		assert.Nil(t, bs)
	})
}
