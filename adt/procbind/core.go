package procbind

import (
	"orglang/go-engine/adt/identity"
	"orglang/go-engine/adt/symbol"
	"orglang/go-engine/adt/uniqref"
	"orglang/go-engine/adt/uniqsym"
	"orglang/go-engine/adt/valkey"
)

// Спецификация связки
type BindSpec struct {
	// channel placeholder (aka variable name)
	ChnlPH symbol.ADT
	// type qualified name (aka variable type)
	TypeQN uniqsym.ADT
}

// Запись связки
type BindRec struct {
	// процесс, в рамках которого связка
	ExecRef uniqref.ADT
	ChnlBS  bindSide
	ChnlPH  symbol.ADT
	ChnlID  identity.ADT
	ExpVK   valkey.ADT
}

type bindSide uint8

const (
	NonSide bindSide = iota
	ProviderSide
	ClientSide
)

func IndexBy[K comparable, V any](getKey func(V) K, vals []V) map[K]V {
	indexed := make(map[K]V)
	for _, val := range vals {
		indexed[getKey(val)] = val
	}
	return indexed
}
