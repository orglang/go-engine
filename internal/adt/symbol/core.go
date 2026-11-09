package symbol

import (
	"orglang/go-engine/internal/adt/contkey"
)

var Zero ADT = ""

type ADT string

func New(str string) ADT {
	if str == "" {
		panic("invalid symbol")
	}
	return ADT(str)
}

func (a ADT) Key() contkey.ADT {
	key := 0
	for _, runeValue := range string(a) {
		key += int(runeValue)
	}
	return contkey.ADT(key)
}
