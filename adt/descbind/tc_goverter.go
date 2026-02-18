package descbind

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
var (
	DataFromRec func(BindRec) (bindRecDS, error)
	DataToRec   func(bindRecDS) (BindRec, error)
)
