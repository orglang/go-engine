package uniqsym

import (
	"hash/fnv"

	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/valkey"
)

type ADT struct {
	sym symbol.ADT
	ns  *ADT
}

func New(sym symbol.ADT) ADT {
	return ADT{sym, nil}
}

func (ns ADT) New(sym symbol.ADT) ADT {
	return ADT{sym, &ns}
}

// symbol
func (a ADT) Sym() symbol.ADT {
	return a.sym
}

// namespace
func (a ADT) NS() ADT {
	if a.ns == nil {
		return empty
	}
	return *a.ns
}

func (a ADT) Key() (valkey.ADT, error) {
	if a == empty {
		panic("invalid value")
	}
	h := fnv.New32a()
	_, err := h.Write([]byte(ConvertToString(a)))
	if err != nil {
		return valkey.Zero, err
	}
	return valkey.ADT(h.Sum32()), nil
}

func (a ADT) Equal(b ADT) bool {
	if a.sym == b.sym && a.ns == b.ns {
		return true
	}
	if a.ns == nil || b.ns == nil {
		return false
	}
	return a.ns.Equal(*b.ns)
}

func (a ADT) String() string {
	return ConvertToString(a)
}

const (
	sepKey = 1
)

var (
	empty ADT
)
