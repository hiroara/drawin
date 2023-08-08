package database_test

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/drawin/database"
)

type entry struct {
	Name string
}

var bucket = []byte("test-bucket")

func openSingleDB(t *testing.T) (*database.SingleDB[*entry], error) {
	path := filepath.Join(t.TempDir(), "test.db")

	return database.OpenSingle[*entry](path, bucket)
}

func TestSingle(t *testing.T) {
	t.Parallel()

	db, err := openDB(t)
	require.NoError(t, err)
	sdb, err := database.Single[*entry](db, bucket)
	require.NoError(t, err)
	assert.NotNil(t, sdb)
	require.NoError(t, sdb.Close())
}

func TestSingleDBView(t *testing.T) {
	t.Parallel()

	sdb, err := openSingleDB(t)
	require.NoError(t, err)
	defer sdb.Close()

	err = sdb.View(func(buc *database.Bucket[*entry]) error {
		assert.NotNil(t, buc)
		return nil
	})
	require.NoError(t, err)
}

func TestSingleDBUpdate(t *testing.T) {
	t.Parallel()

	sdb, err := openSingleDB(t)
	require.NoError(t, err)
	defer sdb.Close()

	err = sdb.Update(func(buc *database.Bucket[*entry]) error {
		assert.NotNil(t, buc)
		return nil
	})
	require.NoError(t, err)
}
