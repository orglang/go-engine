package adapter

import (
	"testing"

	"orglang/go-engine/internal/adt/comproot"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	execeval "orglang/go-engine/internal/proc/execeval/core"
)

func TestExecSnapToGenerated(t *testing.T) {
	id := identity.New()
	got := EvalSnapToGenerated(execeval.EvalSnap{
		CompRef: comproot.SemRef{CompID: id, CompRN: seqnum.ADT(7)},
	})
	if got.CompRef.CompId != identity.ConvertToString(id) || got.CompRef.CompRn != 7 {
		t.Fatalf("unexpected generated snapshot: %#v", got)
	}
}
