package descsem

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
var (
	DataFromBind func(SemBind) (SemBindDS, error)
	DataToBind   func(SemBindDS) (SemBind, error)
)
