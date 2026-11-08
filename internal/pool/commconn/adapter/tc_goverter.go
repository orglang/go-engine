package adapter

import commconn "orglang/go-engine/internal/pool/commconn/core"

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/adt/compsem:Data.*
// goverter:extend orglang/go-engine/internal/pool/commturn/adapter:Data.*
var (
	// goverter:autoMap ConnID
	// goverter:autoMap CompRef
	DataFromRec func(commconn.ConnRec) commconn.ConnRecDS
	// goverter:autoMap ConnID
	DataFromQry func(commconn.ConnQry) commconn.ConnQryDS
	// goverter:autoMap ConnID
	// goverter:ignore Turns
	DataFromMod func(commconn.ConnMod) commconn.ConnModDS
	// goverter:ignore CompID
	// goverter:ignore CompRN
	// goverter:ignore ActNr
	DataToSnap func(commconn.ConnSnapDS) (commconn.ConnSnap, error)
}
