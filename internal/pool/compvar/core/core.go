package core

import (
	"orglang/go-engine/internal/adt/compvar"
	"orglang/go-engine/internal/adt/symbol"
)

type VarSpec = compvar.VarSpec

type VarRec = compvar.VarRec
type StructRec = compvar.StructRec
type LinearRec = compvar.LinearRec

type Mode = compvar.Mode

const (
	StructMode = compvar.StructMode
	LinearMode = compvar.LinearMode
)

const (
	LiabSide  = compvar.LiabSide
	AssetSide = compvar.AssetSide
)

func ConvertRecsToRecMap[T VarRec](recs []T) map[symbol.ADT]T {
	return compvar.ConvertRecsToRecMap(recs)
}
