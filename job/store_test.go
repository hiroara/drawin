package job_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hiroara/drawin/database"
	"github.com/hiroara/drawin/job"
)

func TestStoreCreateJob(t *testing.T) {
	db, err := database.Open([]byte("test-bucket"))
	require.NoError(t, err)
	defer db.Close()

	err = database.Update(db, func(buc *database.Bucket[*job.Job]) error {
		st := job.NewStore(buc)

		url1 := "https://example.com/dir1/image1.jpg"
		url2 := "https://example.com/dir1/image2.jpg"
		url3 := "https://example.com/dir2/image1.jpg" // Same basename with url1

		j, err := st.CreateJob(url1)
		require.NoError(t, err)

		if assert.NotNil(t, j) {
			assert.Equal(t, url1, j.URL)
			assert.Equal(t, "image1.jpg", j.Name)
		}

		j, err = st.CreateJob(url1)
		require.NoError(t, err)
		assert.Nil(t, j)

		j, err = st.CreateJob(url2)
		require.NoError(t, err)

		// Another job created
		if assert.NotNil(t, j) {
			assert.Equal(t, url2, j.URL)
			assert.Equal(t, "image2.jpg", j.Name)
		}

		j, err = st.CreateJob(url3)
		require.NoError(t, err)

		// Another job created with different basename with url1
		if assert.NotNil(t, j) {
			assert.Equal(t, url3, j.URL)
			assert.Equal(t, "image1.1.jpg", j.Name)
		}

		return nil
	})
	require.NoError(t, err)
}

func TestStorePut(t *testing.T) {
	db, err := database.Open([]byte("test-bucket"))
	require.NoError(t, err)
	defer db.Close()

	err = database.Update(db, func(buc *database.Bucket[*job.Job]) error {
		st := job.NewStore(buc)
		require.NoError(t, st.Put(&job.Job{Name: "test.jpg"}))

		val, err := buc.Get([]byte("test.jpg"))
		require.NoError(t, err)
		if assert.NotNil(t, val) {
			assert.Equal(t, "test.jpg", val.Name)
		}
		return nil
	})
	require.NoError(t, err)
}
