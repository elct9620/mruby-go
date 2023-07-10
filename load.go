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
func (s *state) LoadString(code string) (value, error) {
	compiled, err := Compile(bytes.NewBufferString(code))
	if err != nil {
		return nil, err
	}

	return s.Load(bytes.NewBufferString(string(compiled)))
}

// Load execute RITE binary
func (s *state) Load(r io.Reader) (value, error) {
	proc, err := newProc(r)
	if err != nil {
		return nil, err
	}

	return proc.Execute(s)
}

func newProc(r io.Reader) (*proc, error) {
	var header binaryHeader
	err := binary.Read(r, riteByteOrder, &header)
	if err != nil {
		return nil, err
	}

	var executable *irep

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
			executable, err = readIREP(r, header.Size)
		case sectionTypeDebug:
			err = noopSection(r, header.Size)
		case sectionTypeLV:
			err = noopSection(r, header.Size)
		case sectionTypeEOF:
			break
		}

		if err != nil {
			return nil, err
		}

		remain -= header.Size
	}

	return &proc{
		executable: executable,
	}, nil
}

func readIREP(r io.Reader, size uint32) (*irep, error) {
	var riteVersion [4]byte
	err := binary.Read(r, riteByteOrder, &riteVersion)
	if err != nil {
		return nil, err
	}

	irepSize := size - sectionHeaderSize - 4
	binary := make([]byte, irepSize)
	_, err = r.Read(binary)
	if err != nil {
		return nil, err
	}

	sizeStripped := binary[4:]
	return newIREP(bytes.NewBuffer(sizeStripped))
}

func noopSection(r io.Reader, size uint32) error {
	buffer := make([]byte, size-sectionHeaderSize)
	_, err := r.Read(buffer)
	if err != nil {
		return err
	}

	return nil
}
