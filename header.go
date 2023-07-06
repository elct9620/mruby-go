package mruby

import (
	"fmt"
	"unsafe"
)

var binaryHeaderSize = uint32(unsafe.Sizeof(BinaryHeader{}))

type BinaryVersion struct {
	Major [2]byte
	Minor [2]byte
}

func (v BinaryVersion) String() string {
	return fmt.Sprintf("%s.%s", v.Major, v.Minor)
}

type BinaryCompiler struct {
	Name    [4]byte
	Version [4]byte
}

func (c BinaryCompiler) String() string {
	return fmt.Sprintf("%s#%s", c.Name, c.Version)
}

type BinaryHeader struct {
	Identifier [4]byte
	Version    BinaryVersion
	Size       uint32
	Compiler   BinaryCompiler
}

func (h BinaryHeader) String() string {
	return fmt.Sprintf(
		`#<RITE id="%s" version="%s" size="%d" compiler="%s">`,
		h.Identifier,
		h.Version,
		h.Size,
		h.Compiler,
	)
}
