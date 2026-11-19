package adapter

import (
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	webapi "orglang/go-engine/internal/app/web/api"
	compexec "orglang/go-engine/internal/proc/compexec/core"
)

func ExecSnapToGenerated(snap compexec.ExecSnap) webapi.ProcCompSnap {
	return webapi.ProcCompSnap{
		CompRef: webapi.AnyCompRef{
			CompId: identity.ConvertToString(snap.CompRef.CompID),
			CompRn: seqnum.ConvertToInt(snap.CompRef.CompRN),
		},
	}
}
