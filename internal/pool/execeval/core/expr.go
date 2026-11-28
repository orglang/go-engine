package core

import (
	"fmt"
	"log/slog"

	"orglang/go-engine/internal/adt/execroot"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
	commact "orglang/go-engine/internal/pool/commact/core"
)

// ExprKind identifies the kind of a persisted pool term expression.
type ExprKind int

const (
	UnkKind ExprKind = iota
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

type ExprSpec interface {
	IsPoolExpr()
}

// запрос доступа (подписка) со стороны клиента (спрос на доступ)
type AcquireSpec struct {
	CommChnlPH symbol.ADT
	ContExpr   ExprSpec
}

func (s AcquireSpec) spec()     {}
func (AcquireSpec) IsPoolExpr() {}

func (s AcquireSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// предложение доступа (публикация) со стороны провайдера (предложение доступа)
type AcceptSpec struct {
	CommChnlPH symbol.ADT
	ContExpr   ExprSpec
}

func (s AcceptSpec) spec()     {}
func (AcceptSpec) IsPoolExpr() {}

func (s AcceptSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// запрос компетенции (подписка) со стороны нанимателя (спрос на рабочую силу)
type HireSpec struct {
	CommChnlPH symbol.ADT
	// запрашиваемая компетенция
	ProcTermQN uniqsym.ADT
	ContExpr   ExprSpec
}

func (s HireSpec) spec()     {}
func (HireSpec) IsPoolExpr() {}

func (s HireSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// предложение компетенции (публикация) со стороны соискателя (предложение рабочей силы)
type ApplySpec struct {
	CommChnlPH symbol.ADT
	// предлагаемая компетенция
	ProcTermQN uniqsym.ADT
	ContExpr   ExprSpec
}

func (s ApplySpec) spec()     {}
func (ApplySpec) IsPoolExpr() {}

func (s ApplySpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// освобождение доступа (публикация) по инициативе клиента
type ReleaseSpec struct {
	CommChnlPH symbol.ADT
}

func (s ReleaseSpec) spec()     {}
func (ReleaseSpec) IsPoolExpr() {}

func (s ReleaseSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// лишение доступа (подписка) по инициативе провайдера
type DetachSpec struct {
	CommChnlPH symbol.ADT
}

func (s DetachSpec) spec()     {}
func (DetachSpec) IsPoolExpr() {}

func (s DetachSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// лишение должности по инициативе работодателя
type FireSpec struct {
	CommChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
}

func (s FireSpec) spec()     {}
func (FireSpec) IsPoolExpr() {}

func (s FireSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// освобождение должности по инициативе сотрудника
type QuitSpec struct {
	CommChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
}

func (s QuitSpec) spec()     {}
func (QuitSpec) IsPoolExpr() {}

func (s QuitSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

type SpawnSpec struct {
	// ссылка на описание порождаемого процесса
	ProcTermQN uniqsym.ADT
	// ссылки на воплощения потребляемых процессов
	ProcExecQNs []uniqsym.ADT
}

func (s SpawnSpec) spec()     {}
func (SpawnSpec) IsPoolExpr() {}

func (s SpawnSpec) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

type SpawnSpec2 struct {
	// ссылка на описание порождаемого процесса
	ProcTermID identity.ADT
	ProcTermRN seqnum.ADT
	// ссылки на воплощения потребляемых процессов
	ProcExecRefs []execroot.RootRef
}

func (s SpawnSpec2) spec()     {}
func (SpawnSpec2) IsPoolExpr() {}

func (s SpawnSpec2) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("%T%+v", s, s))
}

// ExprRec is the expression-record contract carried by a communication action.
type ExprRec = commact.ExprRec

type AcquireRec struct {
	ContChnlPH symbol.ADT
	ContExpr   ExprSpec
}

func (r AcquireRec) rec()     {}
func (AcquireRec) IsExprRec() {}

type AcceptRec struct {
	ContChnlPH symbol.ADT
	ContExpr   ExprSpec
}

func (r AcceptRec) rec()     {}
func (AcceptRec) IsExprRec() {}

type HireRec struct {
	ContChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
	ContExpr   ExprSpec
}

func (r HireRec) rec()     {}
func (HireRec) IsExprRec() {}

type ApplyRec struct {
	ContChnlPH symbol.ADT
	ProcTermQN uniqsym.ADT
	ContExpr   ExprSpec
}

func (r ApplyRec) rec()     {}
func (ApplyRec) IsExprRec() {}

type ReleaseRec struct {
	ContChnlPH symbol.ADT
}

func (r ReleaseRec) rec()     {}
func (ReleaseRec) IsExprRec() {}

type DetachRec struct {
	ContChnlPH symbol.ADT
}

func (r DetachRec) rec()     {}
func (DetachRec) IsExprRec() {}

type FireRec struct {
	CommChnlPH symbol.ADT
}

func (r FireRec) rec()     {}
func (FireRec) IsExprRec() {}

type QuitRec struct {
	CommChnlPH symbol.ADT
}

func (r QuitRec) rec()     {}
func (QuitRec) IsExprRec() {}

type SpawnRec struct {
	CommChnlPH symbol.ADT
}

func (r SpawnRec) rec()     {}
func (SpawnRec) IsExprRec() {}

func ErrSpecTypeUnexpected(got ExprSpec) error {
	return fmt.Errorf("expr spec unexpected: %T%+v", got, got)
}

func ErrRecTypeUnexpected(got ExprRec) error {
	return fmt.Errorf("expr rec unexpected: %T%+v", got, got)
}

func ErrExprKindUnexpected(got ExprKind) error {
	return fmt.Errorf("expr kind unexpected: %v", got)
}
