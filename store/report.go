package store

import (
	bolt "go.etcd.io/bbolt"

	"github.com/hiroara/drawin/database"
	"github.com/hiroara/drawin/downloader"
	"github.com/hiroara/drawin/marshal"
)

var reportBucketKey = []byte("reports")

func createReportBucket(tx *bolt.Tx) error {
	_, err := tx.CreateBucketIfNotExists(reportBucketKey)
	return err
}

func newReportBucket(tx *bolt.Tx) *database.Bucket[*downloader.Report] {
	return database.NewBucket[*downloader.Report](tx, reportBucketKey, marshal.Msgpack[*downloader.Report]())
}
