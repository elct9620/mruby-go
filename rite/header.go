package rite

import "fmt"

type Version struct {
	Major [2]byte
	Minor [2]byte
}

func (v Version) String() string {
	return fmt.Sprintf("%s.%s", v.Major, v.Minor)
}

type Compiler struct {
	Name    [4]byte
	Version [4]byte
}

func (c Compiler) String() string {
	return fmt.Sprintf("%s#%s", c.Name, c.Version)
}

type Header struct {
	Identifier [4]byte
	Version    Version
	Size       uint32
	Compiler   Compiler
}

func (h Header) String() string {
	return fmt.Sprintf(
		`#<RITE id="%s" version="%s" size="%d" compiler="%s">`,
		h.Identifier,
		h.Version,
		h.Size,
		h.Compiler,
	)
}
