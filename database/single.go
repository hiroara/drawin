package database

import (
	bolt "go.etcd.io/bbolt"
)

type SingleDB[T any] struct {
	*DB
	bucket []byte
}

func OpenSingle[T any](path string, bucket []byte, opts *Options) (*SingleDB[T], error) {
	db, err := Open(path, opts)
	if err != nil {
		return nil, err
	}
	return Single[T](db, bucket)
}

func Single[T any](db *DB, bucket []byte) (*SingleDB[T], error) {
	err := CreateBucket(db, bucket)
	if err != nil {
		return nil, err
	}
	return &SingleDB[T]{DB: db, bucket: bucket}, nil
}

func (db *SingleDB[T]) View(fn func(b *Bucket[T]) error) error {
	return db.DB.View(func(tx *bolt.Tx) error {
		return fn(newBucket[T](tx.Bucket(db.bucket)))
	})
}

func (db *SingleDB[T]) Update(fn func(b *Bucket[T]) error) error {
	return db.DB.Update(func(tx *bolt.Tx) error {
		return fn(newBucket[T](tx.Bucket(db.bucket)))
	})
}
