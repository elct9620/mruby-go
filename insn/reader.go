package insn

import (
	"encoding/binary"
	"io"
)

var DefaultBinaryOrder = binary.BigEndian

type Reader interface {
	io.Reader
	As(any) error
	Uint16() (uint16, error)
	String(int) (string, error)
}

type BinaryReader struct {
	io.Reader
	ByteOrder binary.ByteOrder
}

var _ Reader = &BinaryReader{}

func NewBinaryReader(r io.Reader) *BinaryReader {
	return &BinaryReader{r, DefaultBinaryOrder}
}

func (r *BinaryReader) As(data any) error {
	return binary.Read(r, r.ByteOrder, data)
}

func (r *BinaryReader) Uint16() (uint16, error) {
	var v uint16
	err := r.As(&v)
	return v, err
}

func (r *BinaryReader) String(length int) (string, error) {
	buf := make([]byte, length)
	err := r.As(buf)
	if err != nil {
		return "", err
	}

	return string(buf[0 : length-1]), nil
}
