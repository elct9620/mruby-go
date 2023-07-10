package mruby

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

var riteByteOrder = binary.BigEndian

var (
	ErrSectionOverSize = errors.New("section size larger than binary size")
)

// LoadString execute string
func (s *State) LoadString(code string) (Value, error) {
	compiled, err := Compile(bytes.NewBufferString(code))
	if err != nil {
		return nil, err
	}

	return s.Load(bytes.NewBufferString(string(compiled)))
}

// Load execute RITE binary
func (s *State) Load(r io.Reader) (Value, error) {
	var header binaryHeader
	err := binary.Read(r, riteByteOrder, &header)
	if err != nil {
		return nil, err
	}

	remain := header.Size - binaryHeaderSize
	for remain > sectionHeaderSize {
		var header sectionHeader
		err := binary.Read(r, riteByteOrder, &header)
		if err != nil {
			return nil, err
		}

		isOverSize := header.Size > remain
		if isOverSize {
			return nil, ErrSectionOverSize
		}

		switch header.String() {
		case sectionTypeIREP:
			noopSection(r, header.Size)
		case sectionTypeDebug:
			noopSection(r, header.Size)
		case sectionTypeLV:
			noopSection(r, header.Size)
		case sectionTypeEOF:
			break
		}

		remain -= header.Size
	}

	return header.String(), nil
}

func noopSection(r io.Reader, size uint32) error {
	buffer := make([]byte, size-sectionHeaderSize)
	_, err := r.Read(buffer)
	if err != nil {
		return err
	}

	return nil
}
