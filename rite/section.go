package rite

import "unsafe"

var sectionHeaderSize = uint32(unsafe.Sizeof(SectionHeader{}))

var (
	IdentIREP          = []byte("IREP")
	IdentDebug         = []byte("DBG\x00")
	IdentLocalVariable = []byte("LVAR")
	IdentEOF           = []byte("EOF\x00")
)

const (
	TypeIREP SectionType = iota
	TypeDebug
	TypeLocalVariable
	TypeEOF
)

type SectionType = uint8

type SectionHeader struct {
	Identity [4]byte
	Size     uint32
}

type Section struct {
	header SectionHeader
}

func (s *Section) Type() SectionType {
	if s.header.Identity == [4]byte(IdentIREP) {
		return TypeIREP
	}

	if s.header.Identity == [4]byte(IdentDebug) {
		return TypeDebug
	}

	if s.header.Identity == [4]byte(IdentLocalVariable) {
		return TypeLocalVariable
	}

	return TypeEOF
}

func (s *Section) Size() uint32 {
	return s.header.Size
}
