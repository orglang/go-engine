package adapter

import (
	"reflect"
	"testing"

	"orglang/go-engine/internal/adt/symbol"
	termdef "orglang/go-engine/internal/pool/termdef/core"
)

func TestExpDataRoundTripPreservesNestedDefinition(t *testing.T) {
	spec := termdef.ExpAcquireSpec{
		ConnChnlPH: symbol.ADT("request"),
		ContExp: termdef.ExpHireSpec{
			ConnChnlPH: symbol.ADT("worker"),
			ProcTermQN: "example.worker",
			ContExp: termdef.ExpSpawnSpec{
				ProcTermQN: "example.worker",
			},
		},
	}

	wantRec, err := termdef.ConvertSpecToRec(spec)
	if err != nil {
		t.Fatal(err)
	}
	data, err := dataFromExpRec(wantRec)
	if err != nil {
		t.Fatal(err)
	}
	if len(data.States) != 3 {
		t.Fatalf("stored states = %d, want 3", len(data.States))
	}
	states := make(map[int64]ExpState, len(data.States))
	for _, state := range data.States {
		states[state.ExpCK] = state
	}
	root := states[data.ExpCK]
	if root.SupExpCK != 0 {
		t.Errorf("root parent key = %d, want 0", root.SupExpCK)
	}
	middle := states[root.Spec.ContExpCK]
	if middle.SupExpCK != root.ExpCK {
		t.Errorf("middle parent key = %d, want %d", middle.SupExpCK, root.ExpCK)
	}
	leaf := states[middle.Spec.ContExpCK]
	if leaf.SupExpCK != middle.ExpCK {
		t.Errorf("leaf parent key = %d, want %d", leaf.SupExpCK, middle.ExpCK)
	}

	gotRec, err := dataToExpRec(data)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := termdef.ConvertRecToSpec(gotRec), termdef.ConvertRecToSpec(wantRec); !reflect.DeepEqual(got, want) {
		t.Fatalf("round-trip Definition = %#v, want %#v", got, want)
	}
}

func TestDataToExpRecRejectsMissingContinuation(t *testing.T) {
	_, err := dataToExpRec(ExpData{
		ExpCK: 1,
		States: []ExpState{{
			ExpCK: 1,
			K:     AcquireKind,
			Spec: ExpStateSpec{
				ConnChnlPH: "request",
				ContExpCK:  2,
			},
		}},
	})
	if err == nil {
		t.Fatal("dataToExpRec() error = nil, want missing continuation error")
	}
}

func TestDataToExpRecRejectsInvalidTreeStructure(t *testing.T) {
	tests := []struct {
		name   string
		states []ExpState
	}{
		{
			name: "duplicate expression",
			states: []ExpState{
				{ExpCK: 1, K: ReleaseKind, Spec: ExpStateSpec{ConnChnlPH: "request"}},
				{ExpCK: 1, K: ReleaseKind, Spec: ExpStateSpec{ConnChnlPH: "request"}},
			},
		},
		{
			name: "root has parent",
			states: []ExpState{
				{ExpCK: 1, SupExpCK: 2, K: ReleaseKind, Spec: ExpStateSpec{ConnChnlPH: "request"}},
			},
		},
		{
			name: "unreachable expression",
			states: []ExpState{
				{ExpCK: 1, K: ReleaseKind, Spec: ExpStateSpec{ConnChnlPH: "request"}},
				{ExpCK: 2, K: ReleaseKind, Spec: ExpStateSpec{ConnChnlPH: "worker"}},
			},
		},
		{
			name: "child has unexpected parent",
			states: []ExpState{
				{ExpCK: 1, K: AcquireKind, Spec: ExpStateSpec{ConnChnlPH: "request", ContExpCK: 2}},
				{ExpCK: 2, SupExpCK: 3, K: ReleaseKind, Spec: ExpStateSpec{ConnChnlPH: "worker"}},
			},
		},
		{
			name: "unknown expression kind",
			states: []ExpState{
				{ExpCK: 1, K: 99},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := dataToExpRec(ExpData{ExpCK: 1, States: tt.states})
			if err == nil {
				t.Fatal("dataToExpRec() error = nil, want invalid tree error")
			}
		})
	}
}
