package database

import (
	bolt "go.etcd.io/bbolt"

	"github.com/hiroara/drawin/internal/marshal"
)

type SingleDB[T any] struct {
	*DB
	bucket  []byte
	marshal marshal.Spec[T]
}

func OpenSingle[T any](path string, bucket []byte, m marshal.Spec[T], create bool) (*SingleDB[T], error) {
	db, err := Open(path, create)
	if err != nil {
		return nil, err
	}
	return Single(db, bucket, m)
}

func Single[T any](db *DB, bucket []byte, m marshal.Spec[T]) (*SingleDB[T], error) {
	err := CreateBucket(db, bucket)
	if err != nil {
		return nil, err
	}
	return &SingleDB[T]{DB: db, bucket: bucket, marshal: m}, nil
}

func (db *SingleDB[T]) View(fn func(b *Bucket[T]) error) error {
	return db.DB.View(func(tx *bolt.Tx) error {
		return fn(NewBucket(tx, db.bucket, db.marshal))
	})
}

func (db *SingleDB[T]) Update(fn func(b *Bucket[T]) error) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		return fn(NewBucket(tx, db.bucket, db.marshal))
	})
}
