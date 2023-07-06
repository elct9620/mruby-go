package rite

import (
	"encoding/binary"
	"io"
)

var riteOrder = binary.BigEndian

type RITE struct {
	header Header
}

func Load(r io.Reader) (*RITE, error) {
	rite := &RITE{}

	err := binary.Read(r, riteOrder, &rite.header)
	if err != nil {
		return nil, err
	}

	return rite, nil
}

func (r *RITE) Header() Header {
	return r.header
}
