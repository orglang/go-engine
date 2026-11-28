package adapter

import (
	"orglang/go-engine/internal/adt/execroot"

	connact "orglang/go-engine/internal/pool/connact/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/pool/execeval/adapter:Data.*
var (
	// goverter:map ExecID RootID
	// goverter:map ExecRN RootRN
	DataToRootRef func(ActRec) (execroot.RootRef, error)
	// goverter:map . RootRef
	// goverter:map Exp ValExp
	DataToPubRec func(ActRec) (connact.PubRec, error)
	// goverter:map . RootRef
	// goverter:map Exp ContExp
	DataToSubRec func(ActRec) (connact.SubRec, error)
)
