package descsem

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/uniqsym:Convert.*
var (
	DataFromRec func(SemRec) (SemRecDS, error)
	DataToRec   func(SemRecDS) (SemRec, error)
)
