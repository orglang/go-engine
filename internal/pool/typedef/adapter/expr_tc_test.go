package adapter

import (
	"reflect"
	"testing"

	"orglang/go-engine/internal/adt/contkey"
	typedef "orglang/go-engine/internal/pool/typedef/core"
)

func TestDownRecRoundTrip(t *testing.T) {
	t.Parallel()

	want := typedef.DownRec{
		ExprCK:   contkey.Three,
		ContExp: typedef.OneRec{ExprCK: contkey.One},
	}

	dto := dataFromExprRec(want)
	got, err := dataToExprRec(dto)
	if err != nil {
		t.Fatalf("dataToExprRec() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("round trip mismatch: got %#v, want %#v", got, want)
	}
}
