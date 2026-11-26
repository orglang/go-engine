package adapter

import (
	"testing"

	termdef "orglang/go-engine/internal/proc/termdef/core"
)

func TestDataToExpRecRejectsMalformedTree(t *testing.T) {
	tests := []struct {
		name   string
		states []ExpState
	}{
		{
			name: "duplicate expression",
			states: []ExpState{
				{ExpCK: 1, K: CloseKind, Spec: ExpStateSpec{ContChnlPH: "done"}},
				{ExpCK: 1, K: CloseKind, Spec: ExpStateSpec{ContChnlPH: "done"}},
			},
		},
		{
			name: "root has parent",
			states: []ExpState{
				{ExpCK: 1, SupExpCK: 2, K: CloseKind, Spec: ExpStateSpec{ContChnlPH: "done"}},
			},
		},
		{
			name: "unreachable expression",
			states: []ExpState{
				{ExpCK: 1, K: CloseKind, Spec: ExpStateSpec{ContChnlPH: "done"}},
				{ExpCK: 2, K: CloseKind, Spec: ExpStateSpec{ContChnlPH: "other"}},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := dataToExpRec(ExpData{ExpCK: 1, States: tt.states})
			if err == nil {
				t.Fatal("dataToExpRec() error = nil, want malformed tree error")
			}
		})
	}
}

func TestDataToExpRecAllowsSharedCaseContinuation(t *testing.T) {
	rec, err := dataToExpRec(ExpData{
		ExpCK: 1,
		States: []ExpState{
			{
				ExpCK: 1,
				K:     CaseKind,
				Spec: ExpStateSpec{
					ConnChnlPH: "request",
					ContExpCKs: map[string]int64{"accepted": 2, "rejected": 2},
				},
			},
			{
				ExpCK:    2,
				SupExpCK: 1,
				K:        CloseKind,
				Spec:     ExpStateSpec{ContChnlPH: "done"},
			},
		},
	})
	if err != nil {
		t.Fatalf("dataToExpRec() error = %v", err)
	}
	caseRec, ok := rec.(termdef.ExpCaseRec)
	if !ok {
		t.Fatalf("dataToExpRec() = %T, want ExpCaseRec", rec)
	}
	if len(caseRec.ContExps) != 2 {
		t.Fatalf("case branches = %d, want 2", len(caseRec.ContExps))
	}
}
