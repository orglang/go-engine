package adapter

import (
	"testing"

	termdef "orglang/go-engine/internal/proc/termdef/core"
)

func TestDataToExpRecRejectsMalformedTree(t *testing.T) {
	tests := []struct {
		name   string
		states []ExprState
	}{
		{
			name: "missing root",
			states: []ExprState{
				{ExprCK: 2, K: CloseKind, Spec: ExprStateSpec{ContChnlPH: "done"}},
			},
		},
		{
			name: "duplicate expression",
			states: []ExprState{
				{ExprCK: 1, K: CloseKind, Spec: ExprStateSpec{ContChnlPH: "done"}},
				{ExprCK: 1, K: CloseKind, Spec: ExprStateSpec{ContChnlPH: "done"}},
			},
		},
		{
			name: "root has parent",
			states: []ExprState{
				{ExprCK: 1, SupExpCK: 2, K: CloseKind, Spec: ExprStateSpec{ContChnlPH: "done"}},
			},
		},
		{
			name: "unreachable expression",
			states: []ExprState{
				{ExprCK: 1, K: CloseKind, Spec: ExprStateSpec{ContChnlPH: "done"}},
				{ExprCK: 2, K: CloseKind, Spec: ExprStateSpec{ContChnlPH: "other"}},
			},
		},
		{
			name: "child has unexpected parent",
			states: []ExprState{
				{ExprCK: 1, K: WaitKind, Spec: ExprStateSpec{ContChnlPH: "done", ContExpCK: 2}},
				{ExprCK: 2, SupExpCK: 3, K: CloseKind, Spec: ExprStateSpec{ContChnlPH: "done"}},
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
			_, err := dataToExpRec(ExprData{ExprCK: 1, States: tt.states})
			if err == nil {
				t.Fatal("dataToExpRec() error = nil, want malformed tree error")
			}
		})
	}
}

func TestDataToExpRecAllowsSharedCaseContinuation(t *testing.T) {
	rec, err := dataToExpRec(ExprData{
		ExprCK: 1,
		States: []ExprState{
			{
				ExprCK: 1,
				K:     CaseKind,
				Spec: ExprStateSpec{
					CommChnlPH: "request",
					ContExpCKs: map[string]int64{"accepted": 2, "rejected": 2},
				},
			},
			{
				ExprCK:    2,
				SupExpCK: 1,
				K:        CloseKind,
				Spec:     ExprStateSpec{ContChnlPH: "done"},
			},
		},
	})
	if err != nil {
		t.Fatalf("dataToExpRec() error = %v", err)
	}
	caseRec, ok := rec.(termdef.ExprCaseRec)
	if !ok {
		t.Fatalf("dataToExpRec() = %T, want ExprCaseRec", rec)
	}
	if len(caseRec.ContExps) != 2 {
		t.Fatalf("case branches = %d, want 2", len(caseRec.ContExps))
	}
}
