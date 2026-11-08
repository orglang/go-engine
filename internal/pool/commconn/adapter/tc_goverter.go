package adapter

import (
	commconn "orglang/go-engine/internal/pool/commconn/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/pool/commturn/adapter:Data.*
var (
	// goverter:autoMap ConnRef
	DataFromRec func(commconn.ConnRec) commconn.ConnRecDS
	// goverter:autoMap ConnRef
	DataFromQry func(commconn.ConnQry) commconn.ConnQryDS
	// goverter:autoMap ConnRef
	DataFromMod func(commconn.ConnMod) commconn.ConnModDS
	// goverter:map . ConnRef
	DataToSnap func(commconn.ConnSnapDS) (commconn.ConnSnap, error)
)
