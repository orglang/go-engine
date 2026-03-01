package implvar

import (
	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/revnum"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
)

// human-readable specification of implementation variable
// человекочитаемая спецификация переменной воплощения
type VarSpec struct {
	// channel placeholder
	ChnlPH symbol.ADT
	// implementation qualified name
	ImplQN uniqsym.ADT
}

// machine-readable record of implementation variable
// машиночитаемая запись переменной воплощения
type VarRec interface {
	rec()
}

// linear var
type LinearRec struct {
	// воплощение, в рамках которого связка
	ImplRef implsem.SemRef
	ChnlID  identity.ADT
	ChnlPH  symbol.ADT
	ChnlBS  bindSide
	// ссылка на выражение описания (aka текущий тип канала)
	ExpVK valkey.ADT
}

func (r LinearRec) rec() {}

// structural var
type StructRec struct {
	ChnlRef commsem.SemRef
	ChnlPH  symbol.ADT
	ChnlBS  bindSide
	ExpVK   valkey.ADT
}

func (r StructRec) rec() {}

func (r StructRec) Rewind(offset revnum.ADT) StructRec {
	r.ChnlRef.ChnlON = offset
	return r
}

type bindSide int8

const (
	unkSide bindSide = iota
	Provider
	Client
)

func (bs bindSide) Negate() bindSide {
	return -bs
}

type usageMode uint8

const (
	unkMode usageMode = iota
	Structural
	Linear
)

func IndexBy[K comparable, V any](getKey func(V) K, vals []V) map[K]V {
	indexed := make(map[K]V)
	for _, val := range vals {
		indexed[getKey(val)] = val
	}
	return indexed
}
