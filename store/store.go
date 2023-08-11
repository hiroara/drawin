package store

import (
	bolt "go.etcd.io/bbolt"

	"github.com/hiroara/drawin/database"
	"github.com/hiroara/drawin/job"
	"github.com/hiroara/drawin/reporter"
)

type Store struct {
	db *database.DB
}

func New(db *database.DB) *Store {
	return &Store{db: db}
}

func Open(path string, dbOpts *database.Options) (*Store, error) {
	db, err := database.Open(path, dbOpts)
	if err != nil {
		return nil, err
	}
	return New(db), nil
}

func (s *Store) Initialize() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		if err := createImageBucket(tx); err != nil {
			return err
		}
		if err := createReportBucket(tx); err != nil {
			return err
		}
		return nil
	})
}

func (s *Store) Add(rep *reporter.Report, data []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		return newBucketSet(tx).put(rep, data)
	})
}

func (s *Store) Get(j *job.Job) (*reporter.Report, error) {
	var rep *reporter.Report

	err := s.db.View(func(tx *bolt.Tx) error {
		r, err := newReportBucket(tx).get(j)
		if err != nil {
			return err
		}
		rep = r
		return nil
	})
	if err != nil {
		return nil, err
	}

	return rep, nil
}

func (s *Store) Read(rep *reporter.Report) ([]byte, error) {
	var blob []byte
	err := s.db.View(func(tx *bolt.Tx) error {
		blob = newImageBucket(tx).get(rep)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return blob, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}
