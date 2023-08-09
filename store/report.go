package store

import (
	"github.com/hiroara/carbo/marshal"
	bolt "go.etcd.io/bbolt"

	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/reporter"
)

var reportBucketKey = []byte("reports")

type reportBucket struct {
	*bolt.Bucket

	marshaller marshal.Spec[*reporter.Report]
}

func createReportBucket(tx *bolt.Tx) error {
	_, err := tx.CreateBucketIfNotExists(reportBucketKey)
	return err
}

func newReportBucket(tx *bolt.Tx) *reportBucket {
	return &reportBucket{
		Bucket:     tx.Bucket(reportBucketKey),
		marshaller: marshal.Gob[*reporter.Report](),
	}
}

func (b *reportBucket) get(j *job.Job) (*reporter.Report, error) {
	v := b.Get([]byte(j.URL))
	if v == nil {
		return nil, nil
	}

	return b.marshaller.Unmarshal(v)
}

func (b *reportBucket) put(rep *reporter.Report) error {
	bs, err := b.marshaller.Marshal(rep)
	if err != nil {
		return err
	}

	return b.Put([]byte(rep.URL), bs)
}
