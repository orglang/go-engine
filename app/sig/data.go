package sig

import (
	"smecalculus/rolevod/internal/chnl"
)

type refData struct {
	SigID string `db:"sig_id"`
	Rev   int64  `db:"rev"`
	Title string `db:"title"`
}

type rootData struct {
	SigID string          `db:"sig_id"`
	Title string          `db:"title"`
	Ys    []chnl.SpecData `db:"ys"`
	X     chnl.SpecData   `db:"x"`
	Rev   int64           `db:"rev"`
}

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend smecalculus/rolevod/lib/id:Convert.*
// goverter:extend smecalculus/rolevod/internal/state:Data.*
var (
	DataToRef     func(refData) (Ref, error)
	DataFromRef   func(Ref) refData
	DataToRefs    func([]refData) ([]Ref, error)
	DataFromRefs  func([]Ref) []refData
	DataToRoot    func(rootData) (Root, error)
	DataFromRoot  func(Root) (rootData, error)
	DataToRoots   func([]rootData) ([]Root, error)
	DataFromRoots func([]Root) ([]rootData, error)
)
