package store

import (
	bolt "go.etcd.io/bbolt"

	"github.com/hiroara/drawin/database"
	"github.com/hiroara/drawin/marshal"
	"github.com/hiroara/drawin/reporter"
)

var reportBucketKey = []byte("reports")

func createReportBucket(tx *bolt.Tx) error {
	_, err := tx.CreateBucketIfNotExists(reportBucketKey)
	return err
}

func newReportBucket(tx *bolt.Tx) *database.Bucket[*reporter.Report] {
	return database.NewBucket[*reporter.Report](tx, reportBucketKey, marshal.Msgpack[*reporter.Report]())
}
