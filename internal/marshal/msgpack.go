package marshal

import "github.com/ugorji/go/codec"

type MsgpackSpec[T any] struct {
	handle *codec.MsgpackHandle
}

func Msgpack[T any]() *MsgpackSpec[T] {
	var mh codec.MsgpackHandle
	return &MsgpackSpec[T]{handle: &mh}
}

func (m *MsgpackSpec[T]) Marshal(v T) ([]byte, error) {
	var bs []byte
	enc := codec.NewEncoderBytes(&bs, m.handle)
	if err := enc.Encode(&v); err != nil {
		return nil, err
	}
	return bs, nil
}

func (m *MsgpackSpec[T]) Unmarshal(data []byte) (T, error) {
	dec := codec.NewDecoderBytes(data, m.handle)
	var v T
	if err := dec.Decode(&v); err != nil {
		return v, err
	}
	return v, nil
}
