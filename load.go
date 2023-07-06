package mruby

import (
	"encoding/binary"
	"errors"
	"io"
)

var ErrUnableReadRiteHeader = errors.New("unable to read RITE header")

func ReadRiteHeader(r io.Reader) (header *RiteHeader, err error) {
	header = new(RiteHeader)
	err = binary.Read(r, binary.BigEndian, header)
	if err != nil {
		return nil, ErrUnableReadRiteHeader
	}

	return header, nil
}
