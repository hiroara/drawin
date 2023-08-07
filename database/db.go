package database

import (
	"os"

	bolt "go.etcd.io/bbolt"
)

type DB struct {
	db     *bolt.DB
	file   *os.File
	bucket []byte
}

func Open(bucket []byte) (*DB, error) {
	f, err := os.CreateTemp("", "drawin-*.db")
	if err != nil {
		return nil, err
	}
	db, err := bolt.Open(f.Name(), 0600, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(bucket)
		return err
	})
	if err != nil {
		return nil, err
	}
	return &DB{db: db, file: f, bucket: bucket}, nil
}

func View[T any](db *DB, fn func(b *Bucket[T]) error) error {
	return db.db.View(func(tx *bolt.Tx) error {
		return fn(newBucket[T](tx.Bucket(db.bucket)))
	})
}

func Update[T any](db *DB, fn func(b *Bucket[T]) error) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		return fn(newBucket[T](tx.Bucket(db.bucket)))
	})
}

func (db *DB) Close() error {
	db.db.Close()
	db.file.Close()
	os.Remove(db.file.Name())
	return nil
}
