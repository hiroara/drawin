package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/drawin/database"
)

type entry struct {
	Name string
}

var bucket = []byte("test-bucket")

func TestOpen(t *testing.T) {
	t.Parallel()

	db, err := database.Open(bucket)
	require.NoError(t, err)
	assert.NotNil(t, db)
	require.NoError(t, db.Close())
}

func TestDBRun(t *testing.T) {
	t.Parallel()

	db, err := database.Open(bucket)
	require.NoError(t, err)
	defer db.Close()

	err = database.Update(db, func(buc *database.Bucket[*entry]) error {
		assert.NotNil(t, buc)
		return nil
	})
	require.NoError(t, err)
}
