package adapter

import (
	"testing"

	termdef "orglang/go-engine/internal/proc/termdef/core"
)

func TestDataToExprRecRejectsMalformedTree(t *testing.T) {
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
				{ExprCK: 1, SupExprCK: 2, K: CloseKind, Spec: ExprStateSpec{ContChnlPH: "done"}},
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
				{ExprCK: 1, K: WaitKind, Spec: ExprStateSpec{ContChnlPH: "done", ContExprCK: 2}},
				{ExprCK: 2, SupExprCK: 3, K: CloseKind, Spec: ExprStateSpec{ContChnlPH: "done"}},
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
				t.Fatal("dataToExprRec() error = nil, want malformed tree error")
			}
		})
	}
}

func TestDataToExprRecAllowsSharedCaseContinuation(t *testing.T) {
	rec, err := dataToExprRec(ExprData{
		ExprCK: 1,
		States: []ExprState{
			{
				ExprCK: 1,
				K:      CaseKind,
				Spec: ExprStateSpec{
					CommChnlPH:  "request",
					ContExprCKs: map[string]int64{"accepted": 2, "rejected": 2},
				},
			},
			{
				ExprCK:    2,
				SupExprCK: 1,
				K:         CloseKind,
				Spec:      ExprStateSpec{ContChnlPH: "done"},
			},
		},
	})
	if err != nil {
		t.Fatalf("dataToExprRec() error = %v", err)
	}
	caseRec, ok := rec.(termdef.ExprCaseRec)
	if !ok {
		t.Fatalf("dataToExprRec() = %T, want ExprCaseRec", rec)
	}
	if len(caseRec.ContExprs) != 2 {
		t.Fatalf("case branches = %d, want 2", len(caseRec.ContExprs))
	}
}
