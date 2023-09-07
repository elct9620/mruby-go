package mruby

import (
	"bytes"
	"encoding/binary"
	"io"
)

var DefaultByteOrder = binary.BigEndian

type Reader struct {
	io.Reader
	ByteOrder binary.ByteOrder
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r, DefaultByteOrder}
}

func NewBytesReader(b []byte) *Reader {
	return &Reader{bytes.NewReader(b), DefaultByteOrder}
}

func (r *Reader) ReadAs(data any) error {
	return binary.Read(r, r.ByteOrder, data)
}
