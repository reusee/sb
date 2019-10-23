package sb

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

type Decoder struct {
	r io.Reader
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: r,
	}
}

type DecodeError error

var _ Stream = new(Decoder)

func (d *Decoder) Next() *Token {
	var kind Kind
	if err := binary.Read(d.r, binary.LittleEndian, &kind); errors.Is(err, io.EOF) {
		return nil
	} else if err != nil {
		panic(DecodeError(err))
	}

	var value any
	switch kind {

	case KindBool:
		bs := make([]byte, 1)
		if _, err := io.ReadFull(d.r, bs); err != nil {
			panic(DecodeError(err))
		}
		if bs[0] > 0 {
			value = true
		} else {
			value = false
		}

	case KindInt:
		var i int64
		if err := binary.Read(d.r, binary.LittleEndian, &i); err != nil {
			panic(DecodeError(err))
		}
		value = int(i)

	case KindInt8:
		var i int8
		if err := binary.Read(d.r, binary.LittleEndian, &i); err != nil {
			panic(DecodeError(err))
		}
		value = i

	case KindInt16:
		var i int16
		if err := binary.Read(d.r, binary.LittleEndian, &i); err != nil {
			panic(DecodeError(err))
		}
		value = i

	case KindInt32:
		var i int32
		if err := binary.Read(d.r, binary.LittleEndian, &i); err != nil {
			panic(DecodeError(err))
		}
		value = i

	case KindInt64:
		var i int64
		if err := binary.Read(d.r, binary.LittleEndian, &i); err != nil {
			panic(DecodeError(err))
		}
		value = i

	case KindUint:
		var i uint64
		if err := binary.Read(d.r, binary.LittleEndian, &i); err != nil {
			panic(DecodeError(err))
		}
		value = uint(i)

	case KindUint8:
		var i uint8
		if err := binary.Read(d.r, binary.LittleEndian, &i); err != nil {
			panic(DecodeError(err))
		}
		value = i

	case KindUint16:
		var i uint16
		if err := binary.Read(d.r, binary.LittleEndian, &i); err != nil {
			panic(DecodeError(err))
		}
		value = i

	case KindUint32:
		var i uint32
		if err := binary.Read(d.r, binary.LittleEndian, &i); err != nil {
			panic(DecodeError(err))
		}
		value = i

	case KindUint64:
		var i uint64
		if err := binary.Read(d.r, binary.LittleEndian, &i); err != nil {
			panic(DecodeError(err))
		}
		value = i

	case KindFloat32:
		var i float32
		if err := binary.Read(d.r, binary.LittleEndian, &i); err != nil {
			panic(DecodeError(err))
		}
		value = i

	case KindFloat64:
		var i float64
		if err := binary.Read(d.r, binary.LittleEndian, &i); err != nil {
			panic(DecodeError(err))
		}
		value = i

	case KindString:
		var length uint64
		bs := make([]byte, 1)
		if _, err := io.ReadFull(d.r, bs); err != nil {
			panic(DecodeError(err))
		}
		if bs[0] < 128 {
			length = uint64(bs[0])
		} else {
			bs = make([]byte, ^bs[0])
			if _, err := io.ReadFull(d.r, bs); err != nil {
				panic(DecodeError(err))
			}
			var err error
			length, err = binary.ReadUvarint(bytes.NewReader(bs))
			if err != nil {
				panic(DecodeError(err))
			}
		}
		bs = make([]byte, length)
		if _, err := io.ReadFull(d.r, bs); err != nil {
			panic(DecodeError(err))
		}
		value = string(bs)

	}

	return &Token{
		Kind:  kind,
		Value: value,
	}
}
