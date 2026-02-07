package uniqsym

import (
	"orglang/go-engine/adt/symbol"
)

type ADT struct {
	sym symbol.ADT
	ns  *ADT
}

func New(name symbol.ADT) ADT {
	return ADT{name, nil}
}

func (adt ADT) New(name symbol.ADT) ADT {
	return ADT{name, &adt}
}

// symbol
func (adt ADT) Sym() symbol.ADT {
	return adt.sym
}

// namespace
func (adt ADT) NS() ADT {
	if adt.ns == nil {
		return empty
	}
	return *adt.ns
}

func (adt ADT) Equal(b ADT) bool {
	if adt.sym == b.sym && adt.ns == b.ns {
		return true
	}
	if adt.ns == nil || b.ns == nil {
		return false
	}
	return adt.ns.Equal(*b.ns)
}

var (
	empty ADT
)
