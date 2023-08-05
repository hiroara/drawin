package database

import (
	"os"

	"github.com/hiroara/carbo/marshal"
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

func Update[T any](db *DB, fn func(b *Bucket[T]) error) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		return fn(&Bucket[T]{bucket: tx.Bucket(db.bucket), marshal: marshal.Gob[T]()})
	})
}

func (db *DB) Close() error {
	db.db.Close()
	db.file.Close()
	os.Remove(db.file.Name())
	return nil
}
