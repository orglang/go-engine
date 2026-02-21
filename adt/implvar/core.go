package implvar

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
)

type VarSpec struct {
	// channel placeholder (aka variable name)
	ChnlPH symbol.ADT
	// desc qualified name (aka variable type)
	ImplQN uniqsym.ADT
}

type VarRec struct {
	// процесс, в рамках которого связка
	ImplRef implsem.SemRef
	ChnlBS  bindSide
	ChnlPH  symbol.ADT
	ChnlID  identity.ADT
	ExpVK   valkey.ADT
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
