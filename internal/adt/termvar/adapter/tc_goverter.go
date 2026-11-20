package adapter

import (
	"orglang/go-engine/internal/adt/termvar"
	"orglang/go-engine/internal/adt/typesem"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/symbol:Convert.*
var (
	// goverter:autoMap TypeRef
	DataFromVarRec func(termvar.VarRec) (VarRecDS, error)
	// goverter:map . TypeRef | dataToTypeRef
	DataToVarRec func(VarRecDS) (termvar.VarRec, error)
	// goverter:ignore TypeRN
	dataToTypeRef func(VarRecDS) (typesem.SemRef, error)
)
