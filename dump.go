package mruby

import "fmt"

type RiteHeader struct {
	Identifier [4]byte
	Version    struct {
		Major [2]byte
		Minor [2]byte
	}
	Size     uint32
	Compiler struct {
		Name    [4]byte
		Version [4]byte
	}
}

func (h RiteHeader) String() string {
	return fmt.Sprintf(
		`#<RITE id="%s" version="%s.%s" size="%d" compiler="%s#%s">`,
		h.Identifier,
		h.Version.Major,
		h.Version.Minor,
		h.Size,
		h.Compiler.Name,
		h.Compiler.Version,
	)
}
