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

	var pType poolType
	for i := uint16(0); i < pLen; i++ {
		err = r.ReadAs(&pType)
		if err != nil {
			return err
		}

		reader := poolReaders[pType]
		if reader == nil {
			return ErrUnsupportPoolValueType
		}

		ir.poolValue[i], err = reader(r)
		if err != nil {
			return err
		}

		ir.pLen = uint16(i + 1)
	}

	return nil
}

func poolReadString(r *Reader) (Value, error) {
	sLen, err := r.ReadUint16()
	if err != nil {
		return "", err
	}

	return r.ReadString(int(sLen + 1))
}
