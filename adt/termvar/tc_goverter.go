package termvar

import "orglang/go-engine/adt/typesem"

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/symbol:Convert.*
var (
	// goverter:autoMap TypeRef
	DataFromVarRec func(VarRec) (VarRecDS, error)
	// goverter:map . TypeRef | dataToTypeRef
	DataToVarRec func(VarRecDS) (VarRec, error)
	// goverter:ignore TypeRN
	dataToTypeRef func(VarRecDS) (typesem.SemRef, error)
)
