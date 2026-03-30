package descvar

import "orglang/go-engine/adt/descsem"

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/symbol:Convert.*
var (
	// goverter:autoMap DescRef
	DataFromVarRec func(VarRec) (VarRecDS, error)
	// goverter:map . DescRef | DataToDescRef
	DataToVarRec func(VarRecDS) (VarRec, error)
	// goverter:ignore DescRN
	DataToDescRef func(VarRecDS) (descsem.SemRef, error)
)
