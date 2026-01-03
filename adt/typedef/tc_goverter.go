package typedef

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/orglang/adt/identity:Convert.*
// goverter:extend orglang/orglang/adt/typedef:Convert.*
var (
	ConvertRecToRef  func(DefRec) DefRef
	ConvertSnapToRef func(DefSnap) DefRef
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/orglang/adt/identity:Convert.*
// goverter:extend orglang/orglang/adt/revnum:Convert.*
// goverter:extend orglang/orglang/adt/typedef:Msg.*
// goverter:extend Msg.*
var (
	MsgFromDefSpec  func(DefSpec) DefSpecME
	MsgToDefSpec    func(DefSpecME) (DefSpec, error)
	MsgFromDefRef   func(DefRef) DefRefME
	MsgToDefRef     func(DefRefME) (DefRef, error)
	MsgFromDefRefs  func([]DefRef) []DefRefME
	MsgToDefRefs    func([]DefRefME) ([]DefRef, error)
	MsgFromDefSnap  func(DefSnap) DefSnapME
	MsgToDefSnap    func(DefSnapME) (DefSnap, error)
	MsgFromDefSnaps func([]DefSnap) []DefSnapME
	MsgToDefSnaps   func([]DefSnapME) ([]DefSnap, error)
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/orglang/adt/identity:Convert.*
// goverter:extend orglang/orglang/adt/revnum:Convert.*
// goverter:extend orglang/orglang/adt/typedef:Msg.*
var (
	ViewFromDefRef  func(DefRef) DefRefVP
	ViewToDefRef    func(DefRefVP) (DefRef, error)
	ViewFromDefRefs func([]DefRef) []DefRefVP
	ViewToDefRefs   func([]DefRefVP) ([]DefRef, error)
	ViewFromDefSnap func(DefSnap) DefSnapVP
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/orglang/adt/identity:Convert.*
// goverter:extend orglang/orglang/adt/revnum:Convert.*
// goverter:extend orglang/orglang/adt/typedef:Data.*
// goverter:extend data.*
// goverter:extend DataToTermRef
// goverter:extend DataFromTermRef
var (
	DataToDefRef    func(defRefDS) (DefRef, error)
	DataFromDefRef  func(DefRef) (defRefDS, error)
	DataToDefRefs   func([]defRefDS) ([]DefRef, error)
	DataFromDefRefs func([]DefRef) ([]defRefDS, error)
	DataToDefRec    func(defRecDS) (DefRec, error)
	DataFromDefRec  func(DefRec) (defRecDS, error)
	DataToDefRecs   func([]defRecDS) ([]DefRec, error)
	DataFromDefRecs func([]DefRec) ([]defRecDS, error)
)
