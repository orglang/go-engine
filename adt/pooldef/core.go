package pooldef

import "orglang/orglang/adt/typeexp"

type TermSpec interface {
	poolDef()
}

type CloseSpec struct{}

func (s CloseSpec) poolDef() {}

type WaitSpec struct{}

func (s WaitSpec) poolDef() {}

type SendSpec struct {
	TypeES typeexp.ExpSpec
}

func (s SendSpec) poolDef() {}

type RecvSpec struct {
	TypeES typeexp.ExpSpec
}

func (s RecvSpec) poolDef() {}

type LabSpec struct{}

func (s LabSpec) poolDef() {}

type CaseSpec struct{}

func (s CaseSpec) poolDef() {}

type CallSpec struct{}

func (s CallSpec) poolDef() {}

type SpawnSpec struct{}

func (s SpawnSpec) poolDef() {}

type FwdSpec struct{}

func (s FwdSpec) poolDef() {}
