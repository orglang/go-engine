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

	dto := dataFromExpRec(want)
	got, err := dataToExpRec(dto)
	if err != nil {
		t.Fatalf("dataToExpRec() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("round trip mismatch: got %#v, want %#v", got, want)
	}
}
