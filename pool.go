package mruby

import (
	"errors"
	"io"
)

var ErrUnsupportPoolValueType = errors.New("unsupport pool value type")

type PoolType uint8

const (
	poolTypeString PoolType = iota
	poolTypeInt32
	poolTypeStaticString
	poolTypeInt64
	poolTypeFloat  = 5
	poolTypeBigInt = 7
)

type PoolReader = func(io.Reader) (Value, error)

var poolReaders = map[PoolType]PoolReader{
	poolTypeString: poolReadString,
}

func readPoolValues(ir *irep, r io.Reader) error {
	var pLen uint16
	err := binaryRead(r, &pLen)
	if err != nil {
		return err
	}

	ir.poolValue = make([]Value, pLen)

	var pType PoolType
	for i := 0; i < int(pLen); i++ {
		err = binaryRead(r, &pType)
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

func poolReadString(r io.Reader) (Value, error) {
	var sLen uint16
	err := binaryRead(r, &sLen)
	if err != nil {
		return "", err
	}

	s := make([]byte, sLen+1)
	err = binaryRead(r, s)
	if err != nil {
		return "", err
	}

	return string(s[0:sLen]), nil
}
