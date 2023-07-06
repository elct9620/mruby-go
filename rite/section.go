package rite

import "unsafe"

var sectionHeaderSize = uint32(unsafe.Sizeof(SectionHeader{}))

type SectionHeader struct {
	Identity [4]byte
	Size     uint32
}

type Section struct {
	header SectionHeader
}

func (s *Section) Header() SectionHeader {
	return s.header
}

func (s *Section) Size() uint32 {
	return s.header.Size
}
