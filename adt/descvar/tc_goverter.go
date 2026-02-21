package descvar

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/symbol:Convert.*
var (
	DataFromVarRec func(VarRec) (VarRecDS, error)
)
