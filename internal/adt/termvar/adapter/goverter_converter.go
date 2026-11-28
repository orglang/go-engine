package adapter

import (
	"orglang/go-engine/internal/adt/termvar"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/symbol:Convert.*
var (
	DataFromVarRec func(termvar.VarRec) (VarRec, error)
	DataToVarRec   func(VarRec) (termvar.VarRec, error)
)
