package implvar

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/implsem"
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
type VarRec struct {
	ImplRef implsem.SemRef
	// CommRef commsem.SemRef
	ChnlID identity.ADT
	ChnlPH symbol.ADT
	ChnlBS bindSide

	// Ссылка на выражение описания (aka текущий тип канала).
	//
	// Позитивное значение означает получение.
	// Негативное значение означает лишение.
	// Нулевое значение означает исчерпание.
	ExpVK valkey.ADT
}

type bindSide int8

const (
	unkSide bindSide = iota
	Liab
	Asset
)

func IndexBy[K comparable, V any](getKey func(V) K, vals []V) map[K]V {
	indexed := make(map[K]V)
	for _, val := range vals {
		indexed[getKey(val)] = val
	}
	return indexed
}
