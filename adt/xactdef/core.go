package xactdef

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/uniqref"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/xactexp"
)

type API interface {
	Incept(uniqsym.ADT) (DefRef, error)
	Create(DefSpec) (DefSnap, error)
	Modify(DefSnap) (DefSnap, error)
	RetrieveSnap(DefRef) (DefSnap, error)
	RetreiveRefs() ([]DefRef, error)
}

type DefRef = uniqref.ADT

type DefSpec struct {
	XactQN uniqsym.ADT
	XactES xactexp.ExpSpec
}

type DefRec struct {
	DefRef DefRef
	Title  string
	ExpID  identity.ADT
}

type DefSnap struct {
	DefRef  DefRef
	DefSpec DefSpec
	Title   string
}
