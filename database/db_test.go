package database_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"

	"github.com/hiroara/drawin/database"
)

func openDB(path string, opts *database.Options) (*database.DB, error) {
	return database.Open(path, opts)
}

func TestOpen(t *testing.T) {
	t.Parallel()

	t.Run("NormalCase", func(t *testing.T) {
		t.Parallel()

		db, err := openDB(filepath.Join(t.TempDir(), "test.db"), nil)
		require.NoError(t, err)
		assert.NotNil(t, db)
		require.NoError(t, db.Close())
	})

	t.Run("Create=false,DB=DoesNotExist", func(t *testing.T) {
		t.Parallel()

		_, err := openDB(filepath.Join(t.TempDir(), "test.db"), &database.Options{Create: false})
		require.ErrorIs(t, err, os.ErrNotExist)
	})

	t.Run("Readonly=true,DB=Exists", func(t *testing.T) {
		t.Parallel()

		path := filepath.Join(t.TempDir(), "test.db")

		db, err := openDB(path, nil)
		require.NoError(t, err)
		db.Close()

		db, err = openDB(path, &database.Options{Create: false})
		require.NoError(t, err)
		assert.NotNil(t, db)
		require.NoError(t, db.Close())
	})
}

func TestDBView(t *testing.T) {
	t.Parallel()

	db, err := openDB(filepath.Join(t.TempDir(), "test.db"), nil)
	require.NoError(t, err)
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		assert.NotNil(t, tx)
		return nil
	})
	require.NoError(t, err)
}

func TestDBUpdate(t *testing.T) {
	t.Parallel()

	db, err := openDB(filepath.Join(t.TempDir(), "test.db"), nil)
	require.NoError(t, err)
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		assert.NotNil(t, tx)
		return nil
	})
	require.NoError(t, err)
}
