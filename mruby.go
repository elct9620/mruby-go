package mruby

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

var (
	ErrSectionOverSize = errors.New("section size larger than binary size")
	ErrBinaryEOF       = errors.New("RITE binary is reach end")
)

var riteOrder = binary.BigEndian

type Mrb struct {
	header   BinaryHeader
	sections []*Section
}

func New() *Mrb {
	return &Mrb{}
}

func (mrb *Mrb) LoadString(code string) error {
	compiled, err := Compile(bytes.NewBufferString(code))
	if err != nil {
		return err
	}

	return mrb.Load(bytes.NewBuffer(compiled))
}

func (mrb *Mrb) Load(r io.Reader) error {
	err := binary.Read(r, riteOrder, &mrb.header)
	if err != nil {
		return err
	}

	remain := mrb.Size() - binaryHeaderSize
	for remain > sectionHeaderSize {
		section, err := readSection(r, remain)
		if err != nil {
			return err
		}

		mrb.sections = append(mrb.sections, section)
		remain -= section.Size()
	}

	return nil
}

func (mrb *Mrb) Header() BinaryHeader {
	return mrb.header
}

func (mrb *Mrb) Sections() []*Section {
	return mrb.sections
}

func (mrb *Mrb) Size() uint32 {
	return mrb.header.Size
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
		noopSection(r, section)
		return nil, ErrBinaryEOF
	}

	return section, nil
}

func noopSection(r io.Reader, section *Section) error {
	buffer := make([]byte, section.Size()-sectionHeaderSize)
	_, err := r.Read(buffer)
	if err != nil {
		return err
	}

	return nil
}
