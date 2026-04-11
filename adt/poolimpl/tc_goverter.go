package poolimpl

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/implvar:Convert.*
// goverter:extend orglang/go-engine/adt/implvar:Data.*
var (
	// goverter:map . ImplRef
	DataToImplRec func(implRecDS) (ImplRec, error)
)
