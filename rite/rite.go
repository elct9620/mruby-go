package rite

import (
	"encoding/binary"
	"errors"
	"io"
)

var (
	ErrSectionOverSize = errors.New("section size larger than binary size")
	ErrBinaryEOF       = errors.New("RITE binary is reach end")
)

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

	isOverSize := section.Size() > remain
	if isOverSize {
		return nil, ErrSectionOverSize
	}

	switch section.Type() {
	case TypeIREP:
		noopSection(r, section)
	case TypeDebug:
		noopSection(r, section)
	case TypeLocalVariable:
		noopSection(r, section)
	case TypeEOF:
		return nil, ErrBinaryEOF
	}

	return section, nil
}

func noopSection(r io.Reader, section *Section) error {
	buffer := make([]byte, section.Size())
	_, err := r.Read(buffer)
	if err != nil {
		return err
	}

	return nil
}
