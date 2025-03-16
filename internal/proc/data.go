package proc

import (
	"smecalculus/rolevod/internal/step"
)

type modData struct {
	Locks []lockData
	Bnds  []bndData
	Steps []step.RootData
}

type lockData struct {
	PoolID string
	Rev    int64
}

type bndData struct {
	ProcID  string
	ChnlPH  string
	ChnlID  string
	StateID string
	Rev     int64
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
// goverter:extend smecalculus/rolevod/lib/rev:Convert.*
// goverter:extend smecalculus/rolevod/internal/step:Data.*
var (
	DataFromMod func(Mod) (modData, error)
	DataFromBnd func(Bnd) bndData
)
