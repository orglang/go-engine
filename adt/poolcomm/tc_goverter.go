package poolcomm

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
	DataFromQry func(CommQry) commQryDS
	// goverter:autoMap CommRef
	DataFromMod func(CommMod) commModDS
	// goverter:map . CommRef | DataToSemRef
	DataToRec func(connRecDS) (ConnRec, error)
	// goverter:map . CommRef
	DataToSnap func(commSnapDS) (CommSnap, error)
	// goverter:ignore CommRN
	DataToSemRef func(connRecDS) (commsem.SemRef, error)
)
