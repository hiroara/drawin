package store

import (
	"github.com/hiroara/drawin/reporter"
	bolt "go.etcd.io/bbolt"
)

type bucketSet struct {
	images  *imageBucket
	reports *reportBucket
}

func newBucketSet(tx *bolt.Tx) *bucketSet {
	return &bucketSet{
		images:  newImageBucket(tx),
		reports: newReportBucket(tx),
	}
}

func (bs *bucketSet) put(rep *reporter.Report, data []byte) error {
	if err := bs.images.put(rep, data); err != nil {
		return err
	}
	return bs.reports.put(rep)
}
