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

func (r *Reader) ReadUint16() (uint16, error) {
	var data uint16
	err := r.ReadAs(&data)
	return data, err
}

func (r *Reader) ReadString(length int) (string, error) {
	buffer := make([]byte, length)
	err := r.ReadAs(buffer)
	if err != nil {
		return "", err
	}

	return string(buffer[0 : length-1]), nil
}
