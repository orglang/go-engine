package adapter

import (
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	webapi "orglang/go-engine/internal/app/web/api"
	execeval "orglang/go-engine/internal/proc/execeval/core"
)

func EvalSnapToGenerated(snap execeval.EvalSnap) webapi.ProcCompSnap {
	return webapi.ProcCompSnap{
		CompRef: webapi.AnyCompRef{
			CompId: identity.ConvertToString(snap.CompRef.CompID),
			CompRn: seqnum.ConvertToInt(snap.CompRef.CompRN),
		},
	}
}
