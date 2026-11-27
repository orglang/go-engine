package execvar

import (
	"fmt"

	"orglang/go-engine/internal/adt/comproot"
	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
)

// human-readable specification of execution variable
// человекочитаемая спецификация переменной исполнения
type VarSpec struct {
	// channel placeholder
	ChnlPH symbol.ADT
	// term qualified name
	TermQN uniqsym.ADT
}

// machine-readable record of execution variable
// машиночитаемая запись переменной исполнения
type VarRec interface {
	GetRootRef() comproot.SemRef
	GetConnID() identity.ADT
	GetChnlID() identity.ADT
	GetChnlPH() symbol.ADT
	GetExpCK() contkey.ADT
}

type LinearRec struct {
	RootRef comproot.SemRef
	ConnID  identity.ADT
	ChnlID  identity.ADT
	ChnlPH  symbol.ADT
	ChnlBS  side

	// Ссылка на выражение описания (aka текущий тип канала).
	//
	// Позитивное значение означает получение.
	// Негативное значение означает лишение.
	// Нулевое значение означает исчерпание.
	ExpCK contkey.ADT
}

func (r LinearRec) GetRootRef() comproot.SemRef { return r.RootRef }

func (r LinearRec) GetConnID() identity.ADT { return r.ConnID }

func (r LinearRec) GetChnlID() identity.ADT { return r.ChnlID }

func (r LinearRec) GetChnlPH() symbol.ADT { return r.ChnlPH }

func (r LinearRec) GetExpCK() contkey.ADT { return r.ExpCK }

type StructRec struct {
	RootRef comproot.SemRef
	ConnID  identity.ADT
	ChnlID  identity.ADT
	ChnlPH  symbol.ADT
	ChnlBS  side

	// Ссылка на выражение описания (aka текущий тип канала).
	//
	// Позитивное значение означает получение.
	// Негативное значение означает лишение.
	// Нулевое значение означает исчерпание.
	ExpCK contkey.ADT
}

func (r StructRec) GetRootRef() comproot.SemRef { return r.RootRef }

func (r StructRec) GetConnID() identity.ADT { return r.ConnID }

func (r StructRec) GetChnlID() identity.ADT { return r.ChnlID }

func (r StructRec) GetChnlPH() symbol.ADT { return r.ChnlPH }

func (r StructRec) GetExpCK() contkey.ADT { return r.ExpCK }

type Side int16

type side = Side

const (
	unkSide Side = iota
	LiabSide
	AssetSide
)

type Mode int16

const (
	unkMode Mode = iota
	StructMode
	LinearMode
)

func IndexBy[K comparable, V any](getKey func(V) K, vals []V) map[K]V {
	indexed := make(map[K]V)
	for _, val := range vals {
		indexed[getKey(val)] = val
	}
	return indexed
}

func ErrUnexpectedRecType(got VarRec) error {
	return fmt.Errorf("rec type unexpected: %T", got)
}

func ErrUnexpectedMode(got Mode) error {
	return fmt.Errorf("var mode unexpected: %v", got)
}
