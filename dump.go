package mruby

import (
	"fmt"
	"unsafe"
)

var (
	binaryHeaderSize  = uint32(unsafe.Sizeof(binaryHeader{}))
	sectionHeaderSize = uint32(unsafe.Sizeof(sectionHeader{}))
)

type binaryVersion struct {
	Major [2]byte
	Minor [2]byte
}

func (v binaryVersion) String() string {
	return fmt.Sprintf("%s.%s", v.Major, v.Minor)
}

type binaryCompiler struct {
	Name    [4]byte
	Version [4]byte
}

func (c binaryCompiler) String() string {
	return fmt.Sprintf("%s#%s", c.Name, c.Version)
}

type binaryHeader struct {
	Identifier [4]byte
	Version    binaryVersion
	Size       uint32
	Compiler   binaryCompiler
}

func (h binaryHeader) String() string {
	return fmt.Sprintf(
		`#<RITE id="%s" version="%s" size="%d" compiler="%s">`,
		h.Identifier,
		h.Version,
		h.Size,
		h.Compiler,
	)
}

type sectionType = string

const (
	sectionTypeIREP  sectionType = "IREP"
	sectionTypeDebug             = "DBG\x00"
	sectionTypeLV                = "LVAR"
	sectionTypeEOF               = "EOF\x00"
)

type sectionHeader struct {
	Identity [4]byte
	Size     uint32
}

func (h sectionHeader) String() string {
	return string(h.Identity[:])
}
