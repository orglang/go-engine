package poolconn

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
var (
	// goverter:autoMap CommRef
	DataFromRec func(ConnRec) connRecDS
	// goverter:map . CommRef
	DataToRec func(connRecDS) (ConnRec, error)
)
