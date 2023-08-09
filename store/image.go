package store

import (
	"github.com/hiroara/drawin/reporter"
	bolt "go.etcd.io/bbolt"
)

var imageBucketKey = []byte("images")

type imageBucket struct {
	*bolt.Bucket
}

func createImageBucket(tx *bolt.Tx) error {
	_, err := tx.CreateBucketIfNotExists(imageBucketKey)
	return err
}

func newImageBucket(tx *bolt.Tx) *imageBucket {
	return &imageBucket{Bucket: tx.Bucket(imageBucketKey)}
}

func (b *imageBucket) get(rep *reporter.Report) []byte {
	return b.Get([]byte(rep.URL))
}

func (b *imageBucket) put(rep *reporter.Report, data []byte) error {
	return b.Put([]byte(rep.URL), data)
}
