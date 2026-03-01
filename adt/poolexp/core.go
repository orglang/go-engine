package poolexp

import (
	"fmt"

	"orglang/go-engine/adt/descsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/uniqsym"
)

type ExpSpec interface {
	spec()
}

// запрос доступа со стороны клиента (спрос на доступ)
type AcquireSpec struct {
	CommChnlPH symbol.ADT
}

func (s AcquireSpec) spec() {}

// предложение доступа со стороны провайдера (предложение доступа)
type AcceptSpec struct {
	CommChnlPH symbol.ADT
}

func (s AcceptSpec) spec() {}

// запрос компетенции со стороны нанимателя (спрос на рабочую силу)
type HireSpec struct {
	CommChnlPH symbol.ADT
	// запрашиваемая компетенция
	ProcDescQN uniqsym.ADT
}

func (s HireSpec) spec() {}

// предложение компетенции со стороны соискателя (предложение рабочей силы)
type ApplySpec struct {
	CommChnlPH symbol.ADT
	// предлагаемая компетенция
	ProcDescQN uniqsym.ADT
}

func (s ApplySpec) spec() {}

// освобождение доступа со стороны клиента
type ReleaseSpec struct {
	CommChnlPH symbol.ADT
}

func (s ReleaseSpec) spec() {}

// лишение доступа со стороны провайдера
type DetachSpec struct {
	CommChnlPH symbol.ADT
}

func (s DetachSpec) spec() {}

// лишение должности по инициативе работодателя
type FireSpec struct {
	CommChnlPH symbol.ADT
	ProcDescQN uniqsym.ADT
}

func (s FireSpec) spec() {}

// освобождение должности по инициативе сотрудника
type QuitSpec struct {
	CommChnlPH symbol.ADT
	ProcDescQN uniqsym.ADT
}

func (s QuitSpec) spec() {}

type SpawnSpec struct {
	// ссылка на описание порождаемого процесса
	ProcDescQN uniqsym.ADT
	// ссылки на воплощения потребляемых процессов
	ProcImplQNs []uniqsym.ADT
}

func (s SpawnSpec) spec() {}

type SpawnSpec2 struct {
	// ссылка на описание порождаемого процесса
	ProcDescRef descsem.SemRef
	// ссылки на воплощения потребляемых процессов
	ProcImplRefs []implsem.SemRef
}

func (s SpawnSpec2) spec() {}

type ExpRec interface {
	rec()
}

type AcquireRec struct {
	ContChnlID identity.ADT
}

func (r AcquireRec) rec() {}

type AcceptRec struct {
	ContChnlID identity.ADT
}

func (r AcceptRec) rec() {}

type HireRec struct {
	CommChnlPH symbol.ADT
}

func (r HireRec) rec() {}

type ApplyRec struct {
	CommChnlPH symbol.ADT
}

func (r ApplyRec) rec() {}

type ReleaseRec struct {
	CommChnlPH symbol.ADT
}

func (r ReleaseRec) rec() {}

type DetachRec struct {
	CommChnlPH symbol.ADT
}

func (r DetachRec) rec() {}

type FireRec struct {
	CommChnlPH symbol.ADT
}

func (r FireRec) rec() {}

type QuitRec struct {
	CommChnlPH symbol.ADT
}

func (r QuitRec) rec() {}

type SpawnRec struct {
	CommChnlPH symbol.ADT
}

func (r SpawnRec) rec() {}

func ErrSpecTypeUnexpected(got ExpSpec) error {
	return fmt.Errorf("exp spec unexpected: %T", got)
}

func ErrRecTypeUnexpected(got ExpRec) error {
	return fmt.Errorf("exp rec unexpected: %T", got)
}
