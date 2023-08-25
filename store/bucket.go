package store

import (
	"github.com/hiroara/drawin/database"
	"github.com/hiroara/drawin/downloader"
	bolt "go.etcd.io/bbolt"
)

type bucketSet struct {
	images  *database.Bucket[[]byte]
	reports *database.Bucket[*downloader.Report]
}

func newBucketSet(tx *bolt.Tx) *bucketSet {
	return &bucketSet{
		images:  newImageBucket(tx),
		reports: newReportBucket(tx),
	}
}

func (bs *bucketSet) put(rep *downloader.Report, data []byte) error {
	if data != nil {
		if err := bs.images.Put([]byte(rep.URL), data); err != nil {
			return err
		}
	}
	return bs.reports.Put([]byte(rep.URL), rep)
}
