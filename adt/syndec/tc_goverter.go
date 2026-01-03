package syndec

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/orglang/adt/identity:Convert.*
// goverter:extend orglang/orglang/adt/qualsym:Convert.*
// goverter:extend orglang/orglang/adt/revnum:Convert.*
var (
	DataFromDecRec func(DecRec) (decRecDS, error)
	DataToDecRec   func(decRecDS) (DecRec, error)
)
