package database

import (
	"os"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

type DB struct {
	db *bolt.DB
}

type Options struct {
	Create bool
}

var DefaultOptions = &Options{Create: true}

func Open(path string, opts *Options) (*DB, error) {
	if opts == nil {
		opts = DefaultOptions
	}

	_, err := os.Stat(path)
	if !opts.Create && os.IsNotExist(err) {
		return nil, err
	}

	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return nil, err
	}

	db, err := bolt.Open(path, 0666, nil)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &DB{db: db}, nil
}

func (db *DB) View(fn func(tx *bolt.Tx) error) error {
	return db.db.View(fn)
}

func (db *DB) Update(fn func(tx *bolt.Tx) error) error {
	return db.db.Update(fn)
}

func (db *DB) Close() error {
	db.db.Close()
	return nil
}

func CreateBucket(db *DB, bucket []byte) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		return err
	})
}
