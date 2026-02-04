package poolexp

import (
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/uniqsym"
)

type ExpSpec interface {
	via()
}

type HireSpec struct {
	ProcQN uniqsym.ADT
}

func (s HireSpec) via() {}

type FireSpec struct {
	ProcQN uniqsym.ADT
}

func (s FireSpec) via() {}

type ApplySpec struct {
	ProcQN uniqsym.ADT
}

func (s ApplySpec) via() {}

type QuitSpec struct {
	ProcQN uniqsym.ADT
}

func (s QuitSpec) via() {}

type AcquireSpec struct {
	PoolQN uniqsym.ADT
	BindPH symbol.ADT
}

func (s AcquireSpec) via() {}

type ReleaseSpec struct {
}

func (s ReleaseSpec) via() {}

type AcceptSpec struct {
	PoolQN uniqsym.ADT
	ValPH  symbol.ADT
}

func (s AcceptSpec) via() {}

type DetachSpec struct {
	PoolQN uniqsym.ADT
	ValPH  symbol.ADT
}

func (s DetachSpec) via() {}
