package valkey

import (
	"encoding/binary"
	"hash/fnv"
	"slices"
)

const (
	Zero ADT = iota
	One
	Two
	Three
)

type ADT int64

type Keyable interface {
	Key() ADT
}

func (key ADT) Invert() ADT {
	return -key
}

func Compose(keys ...ADT) (ADT, error) {
	slices.Sort(keys)
	h := fnv.New32a()
	for _, key := range keys {
		err := binary.Write(h, binary.LittleEndian, uint32(key))
		if err != nil {
			return Zero, err
		}
	}
	return ADT(h.Sum32()), nil
}
