package compvar

import (
	"fmt"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/compsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
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
	GetCommRef() commsem.SemRef
	GetChnlID() identity.ADT
	GetChnlPH() symbol.ADT
	GetExpVK() valkey.ADT
}

type LinearRec struct {
	CompRef compsem.SemRef
	CommRef commsem.SemRef
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

func (r LinearRec) GetCommRef() commsem.SemRef { return r.CommRef }

func (r LinearRec) GetChnlID() identity.ADT { return r.ChnlID }

func (r LinearRec) GetChnlPH() symbol.ADT { return r.ChnlPH }

func (r LinearRec) GetExpVK() valkey.ADT { return r.ExpVK }

type StructRec struct {
	CompRef compsem.SemRef
	CommRef commsem.SemRef
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

func (r StructRec) GetCommRef() commsem.SemRef { return r.CommRef }

func (r StructRec) GetChnlID() identity.ADT { return r.ChnlID }

func (r StructRec) GetChnlPH() symbol.ADT { return r.ChnlPH }

func (r StructRec) GetExpVK() valkey.ADT { return r.ExpVK }

type side int16

const (
	unkSide side = iota
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
