package adapter

import (
	"testing"

	"orglang/go-engine/internal/adt/execroot"
	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	execeval "orglang/go-engine/internal/proc/execeval/core"
)

func TestEvalSnapToGenerated(t *testing.T) {
	id := identity.New()
	got := EvalSnapToGenerated(execeval.EvalSnap{
		RootRef: execroot.RootRef{ExecID: id, ExecRN: seqnum.ADT(7)},
	})
	if got.RootRef.ExecId != identity.ConvertToString(id) || got.RootRef.ExecRn != 7 {
		t.Fatalf("unexpected generated snapshot: %#v", got)
	}
}
