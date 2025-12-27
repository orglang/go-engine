package alias

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/orglang/avt/id:Convert.*
// goverter:extend orglang/orglang/avt/rn:Convert.*
// goverter:extend orglang/orglang/avt/sym:Convert.*
var (
	DataFromRoot func(Root) (rootDS, error)
	DataToRoot   func(rootDS) (Root, error)
)
