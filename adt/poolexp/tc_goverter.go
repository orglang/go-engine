package poolexp

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
var (
	DataFromAcquireRec func(AcquireRec) *shiftRecDS
	DataFromAcceptRec  func(AcceptRec) *shiftRecDS

	// goverter:useZeroValueOnPointerInconsistency
	DataToAcquireRec func(*shiftRecDS) (AcquireRec, error)
	// goverter:useZeroValueOnPointerInconsistency
	DataToAcceptRec func(*shiftRecDS) (AcceptRec, error)
)
