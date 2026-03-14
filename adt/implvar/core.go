package implvar

import (
	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/option"
	"orglang/go-engine/adt/seqnum"
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
type VarRec2 interface {
	rec()
}

type VarRec struct {
	ImplRef implsem.SemRef
	CommRef commsem.SemRef
	ChnlID  identity.ADT
	ChnlPH  symbol.ADT
	ChnlBS  bindSide

	// Ссылка на выражение описания (aka текущий тип канала).
	//
	// Позитивное значение означает получение.
	// Негативное значение означает лишение.
	// Нулевое значение означает исчерпание.
	ExpVK valkey.ADT
}

func (r VarRec) Rewind(rn seqnum.ADT) VarRec {
	r.ImplRef.ImplRN = rn
	return r
}

type VarMod struct {
	ChnlID option.ADT[identity.ADT]
	ExpVK  option.ADT[valkey.ADT]
}

type bindSide uint8

const (
	unkSide bindSide = iota
	Provider
	Client
)

func IndexBy[K comparable, V any](getKey func(V) K, vals []V) map[K]V {
	indexed := make(map[K]V)
	for _, val := range vals {
		indexed[getKey(val)] = val
	}
	return indexed
}
