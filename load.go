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

	return s.Load(bytes.NewBufferString(string(compiled)))
}

// Load execute RITE binary
func (s *State) Load(r io.Reader) (Value, error) {
	proc, err := newProc(r)
	if err != nil {
		return nil, err
	}

	return proc.Execute(s)
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
	return newIrep(NewBytesReader(sizeStripped))
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
