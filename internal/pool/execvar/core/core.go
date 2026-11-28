package core

import (
	"orglang/go-engine/internal/adt/execvar"
	"orglang/go-engine/internal/adt/symbol"
)

type VarSpec = execvar.VarSpec

type VarRec = execvar.VarRec
type StructRec = execvar.StructRec
type LinearRec = execvar.LinearRec

type Mode = execvar.Mode

const (
	StructMode = execvar.StructMode
	LinearMode = execvar.LinearMode
)

const (
	LiabSide  = execvar.LiabSide
	AssetSide = execvar.AssetSide
)

func ConvertRecsToRecMap[T VarRec](recs []T) map[symbol.ADT]T {
	return execvar.ConvertRecsToRecMap(recs)
}
