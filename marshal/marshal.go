package marshal

import "github.com/hiroara/carbo/marshal"

type Spec[T any] interface {
	marshal.Spec[T]
}

func Bytes[T marshal.BytesCompatible]() Spec[T] {
	return marshal.Bytes[T]()
}
