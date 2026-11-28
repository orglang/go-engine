package core

import (
	"fmt"
	"log/slog"

	"orglang/go-engine/internal/adt/execroot"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
	connact "orglang/go-engine/internal/pool/connact/core"
)

// ExpKind identifies the kind of a persisted pool term expression.
type ExpKind int

const (
	UnkKind ExpKind = iota
	AcquireKind
	AcceptKind
	HireKind
	ApplyKind
	ReleaseKind
	DetachKind
	FireKind
	QuitKind
	SpawnKind
)

type ExpSpec interface {
	IsPoolExp()
}

// запрос доступа (подписка) со стороны клиента (спрос на доступ)
type AcquireSpec struct {
	ConnChnlPH symbol.ADT
	ContExp    ExpSpec
}

func (s AcquireSpec) spec()    {}
func (AcquireSpec) IsPoolExp() {}

func (s AcquireSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// предложение доступа (публикация) со стороны провайдера (предложение доступа)
type AcceptSpec struct {
	ConnChnlPH symbol.ADT
	ContExp    ExpSpec
}

func (s AcceptSpec) spec()    {}
func (AcceptSpec) IsPoolExp() {}

func (s AcceptSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// запрос компетенции (подписка) со стороны нанимателя (спрос на рабочую силу)
type HireSpec struct {
	ConnChnlPH symbol.ADT
	// запрашиваемая компетенция
	ProcTermQN uniqsym.ADT
	ContExp    ExpSpec
}

func (s HireSpec) spec()    {}
func (HireSpec) IsPoolExp() {}

func (s HireSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// предложение компетенции (публикация) со стороны соискателя (предложение рабочей силы)
type ApplySpec struct {
	ConnChnlPH symbol.ADT
	// предлагаемая компетенция
	ProcTermQN uniqsym.ADT
	ContExp    ExpSpec
}

func (s ApplySpec) spec()    {}
func (ApplySpec) IsPoolExp() {}

func (s ApplySpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// освобождение доступа (публикация) по инициативе клиента
type ReleaseSpec struct {
	ConnChnlPH symbol.ADT
}

func (s ReleaseSpec) spec()    {}
func (ReleaseSpec) IsPoolExp() {}

func (s ReleaseSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// лишение доступа (подписка) по инициативе провайдера
type DetachSpec struct {
	ConnChnlPH symbol.ADT
}

func (s DetachSpec) spec()    {}
func (DetachSpec) IsPoolExp() {}

func (s DetachSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// лишение должности по инициативе работодателя
type FireSpec struct {
	ConnChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
}

func (s FireSpec) spec()    {}
func (FireSpec) IsPoolExp() {}

func (s FireSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// освобождение должности по инициативе сотрудника
type QuitSpec struct {
	ConnChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
}

func (s QuitSpec) spec()    {}
func (QuitSpec) IsPoolExp() {}

func (s QuitSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

type SpawnSpec struct {
	// ссылка на описание порождаемого процесса
	ProcTermQN uniqsym.ADT
	// ссылки на воплощения потребляемых процессов
	ProcCompQNs []uniqsym.ADT
}

func (s SpawnSpec) spec()    {}
func (SpawnSpec) IsPoolExp() {}

func (s SpawnSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

type SpawnSpec2 struct {
	// ссылка на описание порождаемого процесса
	ProcTermID identity.ADT
	ProcTermRN seqnum.ADT
	// ссылки на воплощения потребляемых процессов
	ProcRootRefs []execroot.RootRef
}

func (s SpawnSpec2) spec()    {}
func (SpawnSpec2) IsPoolExp() {}

func (s SpawnSpec2) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// ExpRec is the expression-record contract carried by a connection action.
type ExpRec = connact.ExpRec

type AcquireRec struct {
	ContChnlPH symbol.ADT
	ContExp    ExpSpec
}

func (r AcquireRec) rec()    {}
func (AcquireRec) IsExpRec() {}

type AcceptRec struct {
	ContChnlPH symbol.ADT
	ContExp    ExpSpec
}

func (r AcceptRec) rec()    {}
func (AcceptRec) IsExpRec() {}

type HireRec struct {
	ContChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
	ContExp    ExpSpec
}

func (r HireRec) rec()    {}
func (HireRec) IsExpRec() {}

type ApplyRec struct {
	ContChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
	ContExp    ExpSpec
}

func (r ApplyRec) rec()    {}
func (ApplyRec) IsExpRec() {}

type ReleaseRec struct {
	ContChnlPH symbol.ADT
}

func (r ReleaseRec) rec()    {}
func (ReleaseRec) IsExpRec() {}

type DetachRec struct {
	ContChnlPH symbol.ADT
}

func (r DetachRec) rec()    {}
func (DetachRec) IsExpRec() {}

type FireRec struct {
	ConnChnlPH symbol.ADT
}

func (r FireRec) rec()    {}
func (FireRec) IsExpRec() {}

type QuitRec struct {
	ConnChnlPH symbol.ADT
}

func (r QuitRec) rec()    {}
func (QuitRec) IsExpRec() {}

type SpawnRec struct {
	ConnChnlPH symbol.ADT
}

func (r SpawnRec) rec()    {}
func (SpawnRec) IsExpRec() {}

func ErrSpecTypeUnexpected(got ExpSpec) error {
	return fmt.Errorf("exp spec unexpected: %T%+v", got, got)
}

func ErrRecTypeUnexpected(got ExpRec) error {
	return fmt.Errorf("exp rec unexpected: %T%+v", got, got)
}

func ErrExpKindUnexpected(got ExpKind) error {
	return fmt.Errorf("exp kind unexpected: %v", got)
}
