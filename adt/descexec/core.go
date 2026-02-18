package descexec

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/revnum"
)

type ExecRef struct {
	DescID identity.ADT
	DescRN revnum.ADT
}

func NewExecRef() ExecRef {
	return ExecRef{identity.New(), revnum.New()}
}

type ExecRec struct {
	Ref  ExecRef
	Kind descKind
}

type descKind uint8

const (
	nonDesc descKind = iota
	Xact
	Pool
	Type
	Proc
)
