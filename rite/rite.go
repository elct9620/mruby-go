package rite

import (
	"encoding/binary"
	"errors"
	"io"
)

var ErrSectionOverSize = errors.New("section size larger than binary size")

var riteOrder = binary.BigEndian

type RITE struct {
	header   Header
	sections []*Section
}

func Load(r io.Reader) (*RITE, error) {
	rite := &RITE{}

	err := binary.Read(r, riteOrder, &rite.header)
	if err != nil {
		return nil, err
	}

	remain := rite.header.Size - headerSize
	for remain > sectionHeaderSize {
		section, err := readSection(r, remain)
		if err != nil {
			return nil, err
		}

		rite.sections = append(rite.sections, section)
		remain -= section.Size()
	}

	return rite, nil
}

func (r *RITE) Header() Header {
	return r.header
}

func (r *RITE) Sections() []*Section {
	return r.sections
}

func readSection(r io.Reader, remain uint32) (*Section, error) {
	section := &Section{}
	err := binary.Read(r, riteOrder, &section.header)
	if err != nil {
		return nil, err
	}

	size := section.Size()
	if size > remain {
		return nil, ErrSectionOverSize
	}

	buffer := make([]byte, size)
	_, err = r.Read(buffer)
	if err != nil {
		return section, nil
	}

	return section, nil
}
