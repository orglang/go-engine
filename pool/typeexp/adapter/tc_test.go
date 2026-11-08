package adapter

import (
	"reflect"
	"testing"

	"orglang/go-engine/adt/valkey"
	typeexp "orglang/go-engine/pool/typeexp/core"
)

func TestDownRecRoundTrip(t *testing.T) {
	t.Parallel()

	want := typeexp.DownRec{
		ExpVK:   valkey.Three,
		ContExp: typeexp.OneRec{ExpVK: valkey.One},
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
