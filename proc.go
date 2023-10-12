package mruby

import "io"

type executable interface {
	Execute(state *State) (Value, error)
}

type proc struct {
	executable
}

func newProc(r io.Reader) (*proc, error) {
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
			executable, err = readIRep(r, header.Size)
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
