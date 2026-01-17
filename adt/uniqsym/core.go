package uniqsym

import (
	"orglang/go-runtime/adt/symbol"
)

var Nil ADT

type ADT struct {
	name  symbol.ADT
	space *ADT
}

func New(name symbol.ADT) ADT {
	return ADT{name, &Nil}
}

func (space ADT) New(name symbol.ADT) ADT {
	return ADT{name, &space}
}

func (s ADT) Name() symbol.ADT {
	return s.name
}

func (s ADT) Space() ADT {
	return *s.space
}
