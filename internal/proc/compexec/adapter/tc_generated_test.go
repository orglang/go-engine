package adapter

import (
	"testing"

	"orglang/go-engine/internal/adt/compsem"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	compexec "orglang/go-engine/internal/proc/compexec/core"
)

func TestExecSnapToGenerated(t *testing.T) {
	id := identity.New()
	got := ExecSnapToGenerated(compexec.ExecSnap{
		CompRef: compsem.SemRef{CompID: id, CompRN: seqnum.ADT(7)},
	})
	if got.CompRef.CompId != identity.ConvertToString(id) || got.CompRef.CompRn != 7 {
		t.Fatalf("unexpected generated snapshot: %#v", got)
	}
}
