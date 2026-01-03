package procdef

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend Data.*
// goverter:extend data.*
var (
	DataToTermSpecs   func([]TermSpecDS) ([]TermSpec, error)
	DataFromTermSpecs func([]TermSpec) ([]TermSpecDS, error)
	DataToTermRecs    func([]TermRecDS) ([]TermRec, error)
	DataFromTermRecs  func([]TermRec) ([]TermRecDS, error)
)
