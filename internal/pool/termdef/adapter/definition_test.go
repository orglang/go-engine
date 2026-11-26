package adapter

import (
	"testing"

	"orglang/go-engine/internal/adt/symbol"
	termdef "orglang/go-engine/internal/pool/termdef/core"
)

func TestDefinitionToGeneratedPreservesNestedContinuation(t *testing.T) {
	dto, err := DefinitionToGenerated(termdef.DefSpec{Exp: termdef.ExpAcquireSpec{
		ConnChnlPH: symbol.ADT("request"),
		ContExp:    termdef.ExpReleaseSpec{ConnChnlPH: symbol.ADT("release")},
	}})
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
