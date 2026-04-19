package commexch

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/pool/commturn:Data.*
var (
	// goverter:autoMap CommRef
	DataFromRec func(ExchRec) exchRecDS
	// goverter:autoMap CommRef
	DataFromQry func(ExchQry) exchQryDS
	// goverter:autoMap CommRef
	DataFromMod func(ExchMod) exchModDS
	// goverter:map . CommRef
	DataToSnap func(exchSnapDS) (ExchSnap, error)
)
