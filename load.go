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
func (s *State) LoadString(code string) (Value, error) {
	compiled, err := Compile(bytes.NewBufferString(code))
	if err != nil {
		return nil, err
	}

	return s.LoadIRep(bytes.NewBufferString(string(compiled)))
}

// Load execute RITE binary
func (mrb *State) LoadIRep(r io.Reader) (Value, error) {
	proc, err := readIRep(mrb, r)
	if err != nil {
		return nil, err
	}

	return mrb.TopRun(proc, mrb.TopSelf())
}

func readIRep(mrb *State, r io.Reader) (RProc, error) {
	var header binaryHeader
	err := binaryRead(r, &header)
	if err != nil {
		return nil, err
	}

	var executable *iRep

	remain := header.Size - binaryHeaderSize
	for remain > sectionHeaderSize {
		header, err := readSectionHeader(r, remain)

		switch header.String() {
		case sectionTypeIRep:
			executable, err = readSectionIRep(mrb, r, header.Size)
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

func readSectionHeader(r io.Reader, remain uint32) (*sectionHeader, error) {
	var header sectionHeader
	err := binaryRead(r, &header)
	if err != nil {
		return nil, err
	}

	isOverSize := header.Size > remain
	if isOverSize {
		return nil, ErrSectionOverSize
	}

	return &header, nil
}

func readSectionIRep(mrb *State, r io.Reader, size uint32) (*iRep, error) {
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
	return newIRep(mrb, NewBytesReader(sizeStripped))
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
