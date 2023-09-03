package store

import (
	bolt "go.etcd.io/bbolt"

	"github.com/hiroara/drawin"
	"github.com/hiroara/drawin/internal/database"
	"github.com/hiroara/drawin/marshal"
)

var reportBucketKey = []byte("reports")

func createReportBucket(tx *bolt.Tx) error {
	_, err := tx.CreateBucketIfNotExists(reportBucketKey)
	return err
}

func newReportBucket(tx *bolt.Tx) *database.Bucket[*drawin.Report] {
	return database.NewBucket[*drawin.Report](tx, reportBucketKey, marshal.Msgpack[*drawin.Report]())
}
