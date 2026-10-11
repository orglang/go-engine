package adapter

import (
	termexp "orglang/go-engine/proc/termexp/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend Data.*
// goverter:extend data.*
var (
	DataToExpSpecs   func([]termexp.ExpSpecDS) ([]termexp.ExpSpec, error)
	DataFromExpSpecs func([]termexp.ExpSpec) ([]termexp.ExpSpecDS, error)
	DataToExpRecs    func([]termexp.ExpRecDS) ([]termexp.ExpRec, error)
	DataFromExpRecs  func([]termexp.ExpRec) ([]termexp.ExpRecDS, error)
)
