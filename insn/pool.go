package insn

import (
	"encoding/binary"
	"fmt"
)

type PoolType uint8

const (
	PoolString PoolType = iota
	PoolInt32
	PoolStaticString
	PoolInt64
	PoolFloat  = 5
	PoolBigInt = 7
)

type PoolReader = func(Reader) (any, error)

var poolReaders = map[PoolType]PoolReader{
	PoolString: readStringPool,
	PoolFloat:  readFloatPool,
}

func readPool(r Reader) (any, error) {
	var pType PoolType
	if err := r.As(&pType); err != nil {
		return nil, err
	}

	reader := poolReaders[pType]
	if reader == nil {
		return nil, fmt.Errorf("unsupported pool type: %d", pType)
	}

	return reader(r)
}

func readStringPool(r Reader) (any, error) {
	return r.String()
}

func readFloatPool(r Reader) (any, error) {
	var v float64
	if err := binary.Read(r, binary.LittleEndian, &v); err != nil {
		return nil, err
	}

	return v, nil
}
