package symbol

import (
	"orglang/go-engine/adt/valkey"
)

type ADT string

func New(str string) ADT {
	if str == "" {
		panic("invalid symbol")
	}
	return ADT(str)
}

func (a ADT) Key() valkey.ADT {
	key := 0
	for _, runeValue := range string(a) {
		key += int(runeValue)
	}
	return valkey.ADT(key)
}
