package database

import (
	bolt "go.etcd.io/bbolt"

	"github.com/hiroara/drawin/marshal"
)

type Bucket[T any] struct {
	bucket  *bolt.Bucket
	marshal marshal.Spec[T]
}

func NewBucket[T any](tx *bolt.Tx, bucket []byte, m marshal.Spec[T]) *Bucket[T] {
	return &Bucket[T]{bucket: tx.Bucket(bucket), marshal: m}
}

func (b *Bucket[T]) Get(key []byte) (T, error) {
	bs := b.bucket.Get(key)
	if bs == nil {
		var zero T
		return zero, nil
	}
	return b.marshal.Unmarshal(bs)
}

func (b *Bucket[T]) Put(key []byte, value T) error {
	bs, err := b.marshal.Marshal(value)
	if err != nil {
		return err
	}
	return b.bucket.Put(key, bs)
}
