package adapter

import (
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	webapi "orglang/go-engine/internal/app/web/api"
	execeval "orglang/go-engine/internal/proc/execeval/core"
)

func EvalSnapToGenerated(snap execeval.EvalSnap) webapi.ProcExecSnap {
	return webapi.ProcExecSnap{
		RootRef: webapi.AnyRootRef{
			ExecId: identity.ConvertToString(snap.RootRef.ExecID),
			ExecRn: seqnum.ConvertToInt(snap.RootRef.ExecRN),
		},
	}
}
