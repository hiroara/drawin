package database_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	bolt "go.etcd.io/bbolt"

	"github.com/hiroara/drawin/database"
)

func openDB(t *testing.T) (*database.DB, error) {
	path := filepath.Join(t.TempDir(), "test.db")

	return database.Open(path)
}

func TestOpen(t *testing.T) {
	t.Parallel()

	db, err := openDB(t)
	require.NoError(t, err)
	assert.NotNil(t, db)
	require.NoError(t, db.Close())
}

func TestDBView(t *testing.T) {
	t.Parallel()

	db, err := openDB(t)
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

	path := filepath.Join(t.TempDir(), "test.db")

	db, err := database.Open(path)
	require.NoError(t, err)
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		assert.NotNil(t, tx)
		return nil
	})
	require.NoError(t, err)
}
