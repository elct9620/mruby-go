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
	String() (string, error)
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

func (r *BinaryReader) String() (string, error) {
	sLen, err := r.Uint16()
	if err != nil {
		return "", err
	}

	buf := make([]byte, sLen+1)
	err = r.As(buf)
	if err != nil {
		return "", err
	}

	return string(buf[0:sLen]), nil
}
