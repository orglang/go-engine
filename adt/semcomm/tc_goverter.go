package semcomm

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
var (
	DataFromRef func(CommRef) SemRefDS
	// goverter:autoMap CommRef
	DataFromRec func(SemRec) semRecDS
	// goverter:map . CommRef
	DataToRec func(semRecDS) (SemRec, error)
)
