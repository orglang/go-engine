package syndec

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/adt/revnum:Convert.*
var (
	DataFromDecRec func(DecRec) (decRecDS, error)
	DataToDecRec   func(decRecDS) (DecRec, error)
)
