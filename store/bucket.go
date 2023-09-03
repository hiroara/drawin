package store

import (
	bolt "go.etcd.io/bbolt"

	"github.com/hiroara/drawin"
	"github.com/hiroara/drawin/internal/database"
)

type bucketSet struct {
	images  *database.Bucket[[]byte]
	reports *database.Bucket[*drawin.Report]
}

func newBucketSet(tx *bolt.Tx) *bucketSet {
	return &bucketSet{
		images:  newImageBucket(tx),
		reports: newReportBucket(tx),
	}
}

func (bs *bucketSet) put(rep *drawin.Report, data []byte) error {
	if data != nil {
		if err := bs.images.Put([]byte(rep.URL), data); err != nil {
			return err
		}
	}
	return bs.reports.Put([]byte(rep.URL), rep)
}
