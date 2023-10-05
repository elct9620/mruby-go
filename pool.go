package mruby

import (
	"errors"
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

type PoolReader = func(*Reader) (Value, error)

var poolReaders = map[poolType]PoolReader{
	poolTypeString: poolReadString,
}

func readPoolValues(ir *iRep, r *Reader) error {
	pLen, err := r.ReadUint16()
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

func readPoolValue(r *Reader) (Value, error) {
	var pType poolType
	err := r.ReadAs(&pType)
	if err != nil {
		return nil, err
	}

	reader := poolReaders[pType]
	if reader == nil {
		return nil, ErrUnsupportPoolValueType
	}

	return reader(r)
}

func poolReadString(r *Reader) (Value, error) {
	sLen, err := r.ReadUint16()
	if err != nil {
		return "", err
	}

	return r.ReadString(int(sLen + 1))
}
