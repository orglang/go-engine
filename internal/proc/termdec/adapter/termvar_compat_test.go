package adapter

import (
	"reflect"
	"testing"

	"orglang/go-engine/internal/adt/contkey"
	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/termvar"
	"orglang/go-engine/internal/adt/typesem"
)

func TestTermVarSDKCompatibilityMapping(t *testing.T) {
	want := termvar.VarRec{
		TypeRef: typesem.New(),
		ChnlPH:  symbol.New("channel-1"),
		ExpCK:   contkey.ADT(42),
	}

	dto := MsgFromTermVarRec(want)
	if dto.ExpVK != 42 {
		t.Fatalf("SDK ExpVK = %d, want 42", dto.ExpVK)
	}

	got, err := MsgToTermVarRec(dto)
	if err != nil {
		t.Fatalf("MsgToTermVarRec: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round-trip mismatch:\n got: %#v\nwant: %#v", got, want)
	}
}
