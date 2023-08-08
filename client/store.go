package client

import (
	"github.com/hiroara/carbo/marshal"
	bolt "go.etcd.io/bbolt"

	"github.com/hiroara/drawin/database"
	"github.com/hiroara/drawin/job"
)

type StoreOutput struct {
	db *database.DB
}

var imageBucket = []byte("images")
var reportBucket = []byte("reports")

var jobMarshaller = marshal.Gob[*job.Job]()

func NewStore(db *database.DB) *StoreOutput {
	return &StoreOutput{db: db}
}

func (out *StoreOutput) Add(j *job.Job, data []byte) error {
	bs, err := jobMarshaller.Marshal(j)
	if err != nil {
		return err
	}

	return out.db.Update(func(tx *bolt.Tx) error {
		imgBuc := tx.Bucket(imageBucket)
		if err := imgBuc.Put([]byte(j.Name), data); err != nil {
			return err
		}
		repBuc := tx.Bucket(reportBucket)
		return repBuc.Put([]byte(j.Name), bs)
	})
}

func (out *StoreOutput) Check(j *job.Job) (bool, error) {
	var ok bool
	err := out.db.View(func(tx *bolt.Tx) error {
		repBuc := tx.Bucket(reportBucket)
		bs := repBuc.Get([]byte(j.Name))
		ok = bs != nil
		return nil
	})
	return ok, err
}

func (out *StoreOutput) Prepare() error {
	err := database.CreateBucket(out.db, imageBucket)
	if err != nil {
		return err
	}
	return database.CreateBucket(out.db, reportBucket)
}
