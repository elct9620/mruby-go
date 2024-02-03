package insn

import (
	"encoding/binary"
	"io"
)

var DefaultBinaryOrder = binary.BigEndian

type Reader interface {
	io.Reader
	ReadAs(any) error
	ReadUint16() (uint16, error)
	ReadString(int) (string, error)
}

type BinaryReader struct {
	io.Reader
	ByteOrder binary.ByteOrder
}

var _ Reader = &BinaryReader{}

func NewBinaryReader(r io.Reader) *BinaryReader {
	return &BinaryReader{r, DefaultBinaryOrder}
}

func (r *BinaryReader) ReadAs(data any) error {
	return binary.Read(r, r.ByteOrder, data)
}

func (r *BinaryReader) ReadUint16() (uint16, error) {
	var v uint16
	err := r.ReadAs(&v)
	return v, err
}

func (r *BinaryReader) ReadString(length int) (string, error) {
	buf := make([]byte, length)
	err := r.ReadAs(buf)
	if err != nil {
		return "", err
	}

	return string(buf[0 : length-1]), nil
}
