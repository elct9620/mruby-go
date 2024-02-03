package mruby

import (
	"errors"

	"github.com/elct9620/mruby-go/insn"
)

var ErrUnsupportPoolValueType = errors.New("unsupport pool value type")

type poolType uint8

const (
	poolTypeString poolType = iota
	poolTypeInt32
	poolTypeStaticString
	poolTypeInt64
	poolTypeFloat  = 5
	poolTypeBigInt = 7
)

type PoolReader = func(insn.Reader) (Value, error)

var poolReaders = map[poolType]PoolReader{
	poolTypeString: poolReadString,
}

func readPoolValues(mrb *State, ir *iRep, r insn.Reader) error {
	pLen, err := r.Uint16()
	if err != nil {
		return err
	}

	ir.poolValue = make([]Value, pLen)

	for i := uint16(0); i < pLen; i++ {
		ir.poolValue[i], err = readPoolValue(r)
		if err != nil {
			return err
		}

		ir.pLen = uint16(i + 1)
	}

	return nil
}

func readPoolValue(r insn.Reader) (Value, error) {
	var pType poolType
	err := r.As(&pType)
	if err != nil {
		return nil, err
	}

	reader := poolReaders[pType]
	if reader == nil {
		return nil, ErrUnsupportPoolValueType
	}

	return reader(r)
}

func poolReadString(r insn.Reader) (Value, error) {
	sLen, err := r.Uint16()
	if err != nil {
		return "", err
	}

	return r.String(int(sLen + 1))
}
