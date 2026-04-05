package poolconn

import (
	"orglang/go-engine/adt/commsem"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/adt/identity:Convert.*
// goverter:extend orglang/go-engine/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/adt/poolstep:Data.*
var (
	// goverter:autoMap CommRef
	DataFromRec func(ConnRec) connRecDS
	// goverter:autoMap CommRef
	DataFromQry func(ConnQry) connQryDS
	// goverter:autoMap CommRef
	DataFromMod func(ConnMod) connModDS
	// goverter:map . CommRef | DataToSemRef
	DataToRec func(connRecDS) (ConnRec, error)
	// goverter:map . CommRef
	DataToSnap func(connSnapDS) (ConnSnap, error)
	// goverter:ignore CommRN
	DataToSemRef func(connRecDS) (commsem.SemRef, error)
)
