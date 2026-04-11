package poolexp

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/uniqsym:Convert.*
// goverter:extend DataFromExpSpec
// goverter:extend DataToExpSpec
var (
	DataFromAcquireSpec func(AcquireSpec) *shiftSpecDS
	DataFromAcceptSpec  func(AcceptSpec) *shiftSpecDS
	DataFromHireSpec    func(HireSpec) *laborSpecDS
	DataFromApplySpec   func(ApplySpec) *laborSpecDS

	// goverter:useZeroValueOnPointerInconsistency
	DataToAcquireSpec func(*shiftSpecDS) (AcquireSpec, error)
	// goverter:useZeroValueOnPointerInconsistency
	DataToAcceptSpec func(*shiftSpecDS) (AcceptSpec, error)
	// goverter:useZeroValueOnPointerInconsistency
	DataToHireSpec func(*laborSpecDS) (HireSpec, error)
	// goverter:useZeroValueOnPointerInconsistency
	DataToApplySpec func(*laborSpecDS) (ApplySpec, error)

	DataFromAcquireRec func(AcquireRec) *shiftRecDS
	DataFromAcceptRec  func(AcceptRec) *shiftRecDS

	// goverter:useZeroValueOnPointerInconsistency
	DataToAcquireRec func(*shiftRecDS) (AcquireRec, error)
	// goverter:useZeroValueOnPointerInconsistency
	DataToAcceptRec func(*shiftRecDS) (AcceptRec, error)
)
