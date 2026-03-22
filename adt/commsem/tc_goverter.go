package commsem

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
var (
	// goverter:autoMap CommRef
	DataFromRec func(SemRec) semRecDS
	// goverter:map . CommRef
	DataToRec func(semRecDS) (SemRec, error)
)
