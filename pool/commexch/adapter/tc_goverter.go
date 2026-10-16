package adapter

import (
	commexch "orglang/go-engine/pool/commexch/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/pool/commturn/adapter:Data.*
var (
	// goverter:autoMap CommRef
	DataFromRec func(commexch.ExchRec) commexch.ExchRecDS
	// goverter:autoMap CommRef
	DataFromQry func(commexch.ExchQry) commexch.ExchQryDS
	// goverter:autoMap CommRef
	DataFromMod func(commexch.ExchMod) commexch.ExchModDS
	// goverter:map . CommRef
	DataToSnap func(commexch.ExchSnapDS) (commexch.ExchSnap, error)
)
