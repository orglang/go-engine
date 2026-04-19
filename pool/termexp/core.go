package termexp

import (
	"fmt"
	"log/slog"

	"orglang/go-engine/adt/compsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/termsem"
	"orglang/go-engine/adt/uniqsym"
)

type ExpSpec interface {
	spec()
}

// запрос доступа (подписка) со стороны клиента (спрос на доступ)
type AcquireSpec struct {
	CommChnlPH symbol.ADT
	ContExp    ExpSpec
}

func (s AcquireSpec) spec() {}

func (s AcquireSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// предложение доступа (публикация) со стороны провайдера (предложение доступа)
type AcceptSpec struct {
	CommChnlPH symbol.ADT
	ContExp    ExpSpec
}

func (s AcceptSpec) spec() {}

func (s AcceptSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// запрос компетенции (подписка) со стороны нанимателя (спрос на рабочую силу)
type HireSpec struct {
	CommChnlPH symbol.ADT
	// запрашиваемая компетенция
	ProcTermQN uniqsym.ADT
	ContExp    ExpSpec
}

func (s HireSpec) spec() {}

func (s HireSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// предложение компетенции (публикация) со стороны соискателя (предложение рабочей силы)
type ApplySpec struct {
	CommChnlPH symbol.ADT
	// предлагаемая компетенция
	ProcTermQN uniqsym.ADT
	ContExp    ExpSpec
}

func (s ApplySpec) spec() {}

func (s ApplySpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// освобождение доступа (публикация) по инициативе клиента
type ReleaseSpec struct {
	CommChnlPH symbol.ADT
}

func (s ReleaseSpec) spec() {}

func (s ReleaseSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// лишение доступа (подписка) по инициативе провайдера
type DetachSpec struct {
	CommChnlPH symbol.ADT
}

func (s DetachSpec) spec() {}

func (s DetachSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// лишение должности по инициативе работодателя
type FireSpec struct {
	CommChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
}

func (s FireSpec) spec() {}

func (s FireSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// освобождение должности по инициативе сотрудника
type QuitSpec struct {
	CommChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
}

func (s QuitSpec) spec() {}

func (s QuitSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

type SpawnSpec struct {
	// ссылка на описание порождаемого процесса
	ProcTermQN uniqsym.ADT
	// ссылки на воплощения потребляемых процессов
	ProcCompQNs []uniqsym.ADT
}

func (s SpawnSpec) spec() {}

func (s SpawnSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

type SpawnSpec2 struct {
	// ссылка на описание порождаемого процесса
	ProcTermRef termsem.SemRef
	// ссылки на воплощения потребляемых процессов
	ProcCompRefs []compsem.SemRef
}

func (s SpawnSpec2) spec() {}

func (s SpawnSpec2) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

type ExpRec interface {
	rec()
}

type AcquireRec struct {
	ContChnlID identity.ADT
	ContExp    ExpSpec
}

func (r AcquireRec) rec() {}

type AcceptRec struct {
	ContChnlID identity.ADT
	ContExp    ExpSpec
}

func (r AcceptRec) rec() {}

type HireRec struct {
	ContChnlID identity.ADT
	ProcTermQN uniqsym.ADT
	ContExp    ExpSpec
}

func (r HireRec) rec() {}

type ApplyRec struct {
	ContChnlID identity.ADT
	ProcTermQN uniqsym.ADT
	ContExp    ExpSpec
}

func (r ApplyRec) rec() {}

type ReleaseRec struct {
	ContChnlID identity.ADT
}

func (r ReleaseRec) rec() {}

type DetachRec struct {
	ContChnlID identity.ADT
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
	return fmt.Errorf("exp spec unexpected: %T%+v", got, got)
}

func ErrRecTypeUnexpected(got ExpRec) error {
	return fmt.Errorf("exp rec unexpected: %T%+v", got, got)
}

func ErrExpKindUnexpected(got expKind) error {
	return fmt.Errorf("exp kind unexpected: %v", got)
}
