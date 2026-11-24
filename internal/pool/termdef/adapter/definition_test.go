package adapter

import (
	"testing"

	"orglang/go-engine/internal/adt/symbol"
	termexp "orglang/go-engine/internal/pool/termexp/core"
)

func TestDefinitionToGeneratedPreservesNestedContinuation(t *testing.T) {
	dto, err := DefinitionToGenerated(termexp.DescAcquireSpec{
		ConnChnlPH: symbol.ADT("request"),
		ContExp: termexp.DescReleaseSpec{ConnChnlPH: symbol.ADT("release")},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got, err := dto.Discriminator(); err != nil || got != "acquire" {
		t.Fatalf("unexpected root variant: got %q, err %v", got, err)
	}
	acquire, err := dto.AsPoolTermDefAcquireSpec()
	if err != nil {
		t.Fatal(err)
	}
	if got, err := acquire.ContExpr.Discriminator(); err != nil || got != "release" {
		t.Fatalf("unexpected continuation variant: got %q, err %v", got, err)
	}
}
