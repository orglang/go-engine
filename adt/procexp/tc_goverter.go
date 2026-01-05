package procexp

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend Data.*
// goverter:extend data.*
var (
	DataToExpSpecs   func([]ExpSpecDS) ([]ExpSpec, error)
	DataFromExpSpecs func([]ExpSpec) ([]ExpSpecDS, error)
	DataToExpRecs    func([]ExpRecDS) ([]ExpRec, error)
	DataFromExpRecs  func([]ExpRec) ([]ExpRecDS, error)
)
