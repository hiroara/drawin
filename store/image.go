package store

import (
	bolt "go.etcd.io/bbolt"

	"github.com/hiroara/drawin/internal/database"
	"github.com/hiroara/drawin/marshal"
)

var imageBucketKey = []byte("images")

func createImageBucket(tx *bolt.Tx) error {
	_, err := tx.CreateBucketIfNotExists(imageBucketKey)
	return err
}

func newImageBucket(tx *bolt.Tx) *database.Bucket[[]byte] {
	return database.NewBucket(tx, imageBucketKey, marshal.Bytes[[]byte]())
}
