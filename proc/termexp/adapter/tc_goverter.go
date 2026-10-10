package adapter

import (
	proctermexp "orglang/go-engine/proc/termexp/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend Data.*
// goverter:extend data.*
var (
	DataToExpSpecs   func([]proctermexp.ExpSpecDS) ([]proctermexp.ExpSpec, error)
	DataFromExpSpecs func([]proctermexp.ExpSpec) ([]proctermexp.ExpSpecDS, error)
	DataToExpRecs    func([]proctermexp.ExpRecDS) ([]proctermexp.ExpRec, error)
	DataFromExpRecs  func([]proctermexp.ExpRec) ([]proctermexp.ExpRecDS, error)
)
