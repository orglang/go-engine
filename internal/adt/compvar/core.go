package compvar

import (
	"fmt"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
	"orglang/go-engine/internal/adt/valkey"
)

// human-readable specification of computation variable
// человекочитаемая спецификация переменной вычисления
type VarSpec struct {
	// channel placeholder
	ChnlPH symbol.ADT
	// term qualified name
	TermQN uniqsym.ADT
}

// machine-readable record of computation variable
// машиночитаемая запись переменной вычисления
type VarRec interface {
	GetConnID() identity.ADT
	GetChnlID() identity.ADT
	GetChnlPH() symbol.ADT
	GetExpVK() valkey.ADT
}

type LinearRec struct {
	CompRef compsem.SemRef
	ConnID identity.ADT
	ChnlID  identity.ADT
	ChnlPH  symbol.ADT
	ChnlBS  side

	// Ссылка на выражение описания (aka текущий тип канала).
	//
	// Позитивное значение означает получение.
	// Негативное значение означает лишение.
	// Нулевое значение означает исчерпание.
	ExpVK valkey.ADT
}

func (r LinearRec) GetConnID() identity.ADT { return r.ConnID }

func (r LinearRec) GetChnlID() identity.ADT { return r.ChnlID }

func (r LinearRec) GetChnlPH() symbol.ADT { return r.ChnlPH }

func (r LinearRec) GetExpVK() valkey.ADT { return r.ExpVK }

type StructRec struct {
	CompRef compsem.SemRef
	ConnID identity.ADT
	ChnlID  identity.ADT
	ChnlPH  symbol.ADT
	ChnlBS  side

	// Ссылка на выражение описания (aka текущий тип канала).
	//
	// Позитивное значение означает получение.
	// Негативное значение означает лишение.
	// Нулевое значение означает исчерпание.
	ExpVK valkey.ADT
}

func (r StructRec) GetConnID() identity.ADT { return r.ConnID }

func (r StructRec) GetChnlID() identity.ADT { return r.ChnlID }

func (r StructRec) GetChnlPH() symbol.ADT { return r.ChnlPH }

func (r StructRec) GetExpVK() valkey.ADT { return r.ExpVK }

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
