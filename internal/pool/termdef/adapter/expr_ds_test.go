package adapter

import (
	"reflect"
	"testing"

	"orglang/go-engine/internal/adt/symbol"
	"orglang/go-engine/internal/adt/uniqsym"
	termdef "orglang/go-engine/internal/pool/termdef/core"
)

func TestExprDataRoundTripPreservesNestedDefinition(t *testing.T) {
	workerQN, err := uniqsym.ConvertFromString("example.worker")
	if err != nil {
		t.Fatal(err)
	}
	spec := termdef.ExprAcquireSpec{
		CommChnlPH: symbol.ADT("request"),
		ContExp: termdef.ExprHireSpec{
			CommChnlPH: symbol.ADT("worker"),
			ProcTermQN: workerQN,
			ContExp: termdef.ExprSpawnSpec{
				ProcTermQN: workerQN,
			},
		},
	}

	wantRec, err := termdef.ConvertSpecToRec(spec)
	if err != nil {
		t.Fatal(err)
	}
	data, err := dataFromExprRec(wantRec)
	if err != nil {
		t.Fatal(err)
	}
	if len(data.States) != 3 {
		t.Fatalf("stored states = %d, want 3", len(data.States))
	}
	states := make(map[int64]ExprState, len(data.States))
	for _, state := range data.States {
		states[state.ExprCK] = state
	}
	root := states[data.ExprCK]
	if root.SupExprCK != 0 {
		t.Errorf("root parent key = %d, want 0", root.SupExprCK)
	}
	middle := states[root.Spec.ContExprCK]
	if middle.SupExprCK != root.ExprCK {
		t.Errorf("middle parent key = %d, want %d", middle.SupExprCK, root.ExprCK)
	}
	leaf := states[middle.Spec.ContExprCK]
	if leaf.SupExprCK != middle.ExprCK {
		t.Errorf("leaf parent key = %d, want %d", leaf.SupExprCK, middle.ExprCK)
	}

	gotRec, err := dataToExprRec(data)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := termdef.ConvertRecToSpec(gotRec), termdef.ConvertRecToSpec(wantRec); !reflect.DeepEqual(got, want) {
		t.Fatalf("round-trip Definition = %#v, want %#v", got, want)
	}
}

func TestDataToExprRecRejectsMissingContinuation(t *testing.T) {
	_, err := dataToExprRec(ExprData{
		ExprCK: 1,
		States: []ExprState{{
			ExprCK: 1,
			K:     AcquireKind,
			Spec: ExprStateSpec{
				CommChnlPH: "request",
				ContExprCK:  2,
			},
		}},
	})
	if err == nil {
		t.Fatal("dataToExprRec() error = nil, want missing continuation error")
	}
}

func TestDataToExprRecRejectsInvalidTreeStructure(t *testing.T) {
	tests := []struct {
		name   string
		states []ExprState
	}{
		{
			name: "missing root",
			states: []ExprState{
				{ExprCK: 2, K: ReleaseKind, Spec: ExprStateSpec{CommChnlPH: "request"}},
			},
		},
		{
			name: "duplicate expression",
			states: []ExprState{
				{ExprCK: 1, K: ReleaseKind, Spec: ExprStateSpec{CommChnlPH: "request"}},
				{ExprCK: 1, K: ReleaseKind, Spec: ExprStateSpec{CommChnlPH: "request"}},
			},
		},
		{
			name: "root has parent",
			states: []ExprState{
				{ExprCK: 1, SupExprCK: 2, K: ReleaseKind, Spec: ExprStateSpec{CommChnlPH: "request"}},
			},
		},
		{
			name: "unreachable expression",
			states: []ExprState{
				{ExprCK: 1, K: ReleaseKind, Spec: ExprStateSpec{CommChnlPH: "request"}},
				{ExprCK: 2, K: ReleaseKind, Spec: ExprStateSpec{CommChnlPH: "worker"}},
			},
		},
		{
			name: "child has unexpected parent",
			states: []ExprState{
				{ExprCK: 1, K: AcquireKind, Spec: ExprStateSpec{CommChnlPH: "request", ContExprCK: 2}},
				{ExprCK: 2, SupExprCK: 3, K: ReleaseKind, Spec: ExprStateSpec{CommChnlPH: "worker"}},
			},
		},
		{
			name: "unknown expression kind",
			states: []ExprState{
				{ExprCK: 1, K: 99},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := dataToExprRec(ExprData{ExprCK: 1, States: tt.states})
			if err == nil {
				t.Fatal("dataToExprRec() error = nil, want invalid tree error")
			}
		})
	}
}
