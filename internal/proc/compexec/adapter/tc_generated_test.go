package adapter

import (
	"testing"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/compsem"
	compexec "orglang/go-engine/internal/proc/compexec/core"
)

func TestExecSnapToGenerated(t *testing.T) {
	id, err := identity.ConvertFromString("execution-1")
	if err != nil {
		t.Fatal(err)
	}
	got := ExecSnapToGenerated(compexec.ExecSnap{
		CompRef: compsem.SemRef{CompID: id, CompRN: seqnum.ADT(7)},
	})
	if got.CompRef.CompId != "execution-1" || got.CompRef.CompRn != 7 {
		t.Fatalf("unexpected generated snapshot: %#v", got)
	}
}
