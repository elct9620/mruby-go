package insn

import "fmt"

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
