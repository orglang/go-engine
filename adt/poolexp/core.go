package poolexp

import (
	"fmt"

	"orglang/go-engine/adt/descsem"
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/uniqsym"
)

type ExpSpec interface {
	spec()
}

type HireSpec struct {
	CommChnlPH symbol.ADT
	ProcDescQN uniqsym.ADT
}

func (s HireSpec) spec() {}

type FireSpec struct {
	CommChnlPH symbol.ADT
	ProcDescQN uniqsym.ADT
}

func (s FireSpec) spec() {}

type ApplySpec struct {
	CommChnlPH symbol.ADT
	ProcDescQN uniqsym.ADT
}

func (s ApplySpec) spec() {}

type QuitSpec struct {
	CommChnlPH symbol.ADT
	ProcDescQN uniqsym.ADT
}

func (s QuitSpec) spec() {}

type AcquireSpec struct {
	PoolQN uniqsym.ADT
	BindPH symbol.ADT
}

func (s AcquireSpec) spec() {}

type ReleaseSpec struct {
}

func (s ReleaseSpec) spec() {}

type AcceptSpec struct {
	PoolQN uniqsym.ADT
	ValPH  symbol.ADT
}

func (s AcceptSpec) spec() {}

type DetachSpec struct {
	PoolQN uniqsym.ADT
	ValPH  symbol.ADT
}

func (s DetachSpec) spec() {}

type SpawnSpec struct {
	// ссылка на описание порождаемого процесса
	ProcDescRef descsem.SemRef
	// ссылки на воплощения потребляемых процессов
	ProcImplRefs []implsem.SemRef
}

func (s SpawnSpec) spec() {}

func ErrExpTypeUnexpected(got ExpSpec) error {
	return fmt.Errorf("exp spec unexpected: %T", got)
}
