package uniqref

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/revnum"
)

type ADT struct {
	ID identity.ADT
	RN revnum.ADT
}

func New() ADT {
	return ADT{identity.New(), revnum.New()}
}
