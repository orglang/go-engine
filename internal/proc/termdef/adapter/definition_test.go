package adapter

import (
	"testing"

	"orglang/go-engine/internal/adt/symbol"
	termexp "orglang/go-engine/internal/proc/termexp/core"
)

func TestDefinitionToGeneratedPreservesNestedContinuation(t *testing.T) {
	dto, err := DefinitionToGenerated(termexp.DescWaitSpec{
		ContChnlPH: symbol.ADT("next"),
		ContExp: termexp.DescCloseSpec{ContChnlPH: symbol.ADT("done")},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, err := dto.Discriminator(); err != nil || got != "wait" {
		t.Fatalf("unexpected root variant: got %q, err %v", got, err)
	}
	wait, err := dto.AsProcTermDefWaitSpec()
	if err != nil {
		t.Fatal(err)
	}
	if got, err := wait.ContExpr.Discriminator(); err != nil || got != "close" {
		t.Fatalf("unexpected continuation variant: got %q, err %v", got, err)
	}
}
