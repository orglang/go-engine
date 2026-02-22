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
	// воплощение, в рамках которого связка
	ImplRef implsem.SemRef
	ChnlBS  bindSide
	ChnlPH  symbol.ADT
	ChnlID  identity.ADT
	// ссылка на выражение описания (aka текущее состояние канала)
	ExpVK valkey.ADT
}

type bindSide uint8

const (
	unknown bindSide = iota
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
