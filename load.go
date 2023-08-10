package mruby

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

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
	err := binaryRead(r, &header)
	if err != nil {
		return nil, err
	}

	var executable *irep

	remain := header.Size - binaryHeaderSize
	for remain > sectionHeaderSize {
		var header sectionHeader
		err := binaryRead(r, &header)
		if err != nil {
			return nil, err
		}

		isOverSize := header.Size > remain
		if isOverSize {
			return nil, ErrSectionOverSize
		}

		switch header.String() {
		case sectionTypeIrep:
			executable, err = readIrep(r, header.Size)
		case sectionTypeDebug:
			err = noopSection(r, header.Size)
		case sectionTypeLv:
			err = noopSection(r, header.Size)
		case sectionTypeEof:
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

func readIrep(r io.Reader, size uint32) (*irep, error) {
	var riteVersion [4]byte
	err := binaryRead(r, &riteVersion)
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
	return newIrep(bytes.NewBuffer(sizeStripped))
}

func noopSection(r io.Reader, size uint32) error {
	buffer := make([]byte, size-sectionHeaderSize)
	_, err := r.Read(buffer)
	if err != nil {
		return err
	}

	return nil
}

func binaryRead(r io.Reader, data any) error {
	return binary.Read(r, binary.BigEndian, data)
}
