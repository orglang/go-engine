package synonym

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend orglang/go-engine/adt/revnum:Convert.*
// goverter:extend orglang/go-engine/adt/valkey:Convert.*
var (
	DataFromRec func(Rec) (recDS, error)
	DataToRec   func(recDS) (Rec, error)
)
