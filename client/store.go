package client

import (
	"github.com/hiroara/carbo/marshal"
	bolt "go.etcd.io/bbolt"

	"github.com/hiroara/drawin/database"
	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/reporter"
)

type StoreOutput struct {
	db *database.DB
}

var imageBucket = []byte("images")
var reportBucket = []byte("reports")

var reportMarshaller = marshal.Gob[*reporter.Report]()

func NewStore(db *database.DB) *StoreOutput {
	return &StoreOutput{db: db}
}

func (out *StoreOutput) Add(rep *reporter.Report, data []byte) error {
	bs, err := reportMarshaller.Marshal(rep)
	if err != nil {
		return err
	}

	return out.db.Update(func(tx *bolt.Tx) error {
		imgBuc := tx.Bucket(imageBucket)
		if err := imgBuc.Put([]byte(rep.URL), data); err != nil {
			return err
		}
		repBuc := tx.Bucket(reportBucket)
		return repBuc.Put([]byte(rep.URL), bs)
	})
}

func (out *StoreOutput) Get(j *job.Job) (*reporter.Report, error) {
	var rep *reporter.Report
	err := out.db.View(func(tx *bolt.Tx) error {
		repBuc := tx.Bucket(reportBucket)
		bs := repBuc.Get([]byte(j.URL))
		if bs == nil {
			return nil
		}
		r, err := reportMarshaller.Unmarshal(bs)
		if err != nil {
			return err
		}
		r.Result = reporter.Cached
		rep = r
		return nil
	})
	return rep, err
}

func (out *StoreOutput) Prepare() error {
	err := database.CreateBucket(out.db, imageBucket)
	if err != nil {
		return err
	}
	return database.CreateBucket(out.db, reportBucket)
}
