package adapter

import (
	"orglang/go-engine/internal/adt/execroot"

	commact "orglang/go-engine/internal/pool/commact/core"
)

// goverter:variables
// goverter:output:format assign-variable
// goverter:extend orglang/go-engine/internal/adt/identity:Convert.*
// goverter:extend orglang/go-engine/internal/adt/seqnum:Convert.*
// goverter:extend orglang/go-engine/internal/pool/execeval/adapter:Data.*
var (
	// goverter:map ExecID ExecID
	// goverter:map ExecRN ExecRN
	DataToRootRef func(ActRec) (execroot.RootRef, error)
	// goverter:map . RootRef
	// goverter:map Expr ValExp
	DataToPubRec func(ActRec) (commact.PubRec, error)
	// goverter:map . RootRef
	// goverter:map Expr ContExp
	DataToSubRec func(ActRec) (commact.SubRec, error)
)
