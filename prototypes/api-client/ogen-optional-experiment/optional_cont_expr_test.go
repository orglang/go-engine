package ogenclient

import (
	"encoding/json"
	"testing"
)

// This experiment deliberately weakens the Pool.Term wire contract:
// cont_expr may be omitted, but when present it must contain a valid ExprSpec.
func TestOptionalPoolTermContExpr(t *testing.T) {
	tests := []struct {
		name     string
		newValue func() any
		missing  string
		present  string
	}{
		{
			name: "hire",
			newValue: func() any { return &PoolTermHireSpec{} },
			missing: `{"conn_ph":"c","proc_qn":"p"}`,
			present: `{"conn_ph":"c","proc_qn":"p","cont_expr":{"kind":"fire","conn_ph":"c"}}`,
		},
		{
			name: "apply",
			newValue: func() any { return &PoolTermApplySpec{} },
			missing: `{"conn_ph":"c","proc_qn":"p"}`,
			present: `{"conn_ph":"c","proc_qn":"p","cont_expr":{"kind":"fire","conn_ph":"c"}}`,
		},
		{
			name: "acquire",
			newValue: func() any { return &PoolTermAcquireSpec{} },
			missing: `{"conn_ph":"c"}`,
			present: `{"conn_ph":"c","cont_expr":{"kind":"fire","conn_ph":"c"}}`,
		},
		{
			name: "accept",
			newValue: func() any { return &PoolTermAcceptSpec{} },
			missing: `{"conn_ph":"c"}`,
			present: `{"conn_ph":"c","cont_expr":{"kind":"fire","conn_ph":"c"}}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("missing is accepted", func(t *testing.T) {
				if err := json.Unmarshal([]byte(tt.missing), tt.newValue()); err != nil {
					t.Fatalf("expected missing optional cont_expr to decode: %v", err)
				}
			})
			t.Run("present value round-trips", func(t *testing.T) {
				got := tt.newValue()
				if err := json.Unmarshal([]byte(tt.present), got); err != nil {
					t.Fatalf("decode present cont_expr: %v", err)
				}
				encoded, err := json.Marshal(got)
				if err != nil {
					t.Fatalf("encode present cont_expr: %v", err)
				}
				decoded := tt.newValue()
				if err := json.Unmarshal(encoded, decoded); err != nil {
					t.Fatalf("decode encoded value: %v", err)
				}
			})
			t.Run("null is rejected for non-nullable optional field", func(t *testing.T) {
				var payload map[string]json.RawMessage
				if err := json.Unmarshal([]byte(tt.missing), &payload); err != nil {
					t.Fatal(err)
				}
				payload["cont_expr"] = json.RawMessage("null")
				withNull, err := json.Marshal(payload)
				if err != nil {
					t.Fatal(err)
				}
				if err := json.Unmarshal(withNull, tt.newValue()); err == nil {
					t.Fatal("expected explicit null cont_expr to be rejected")
				}
			})
		})
	}
}

func TestOptionalPoolTermContExprNestedRecursion(t *testing.T) {
	// A recursive optional edge must remain traversable when it is present.
	payload := `{"conn_ph":"c","proc_qn":"p","cont_expr":{"kind":"hire","conn_ph":"c","proc_qn":"nested","cont_expr":{"kind":"fire","conn_ph":"c"}}}`
	got := &PoolTermHireSpec{}
	if err := json.Unmarshal([]byte(payload), got); err != nil {
		t.Fatalf("decode nested recursive cont_expr: %v", err)
	}
	encoded, err := json.Marshal(got)
	if err != nil {
		t.Fatalf("encode nested recursive cont_expr: %v", err)
	}
	decoded := &PoolTermHireSpec{}
	if err := json.Unmarshal(encoded, decoded); err != nil {
		t.Fatalf("decode nested round-trip: %v", err)
	}
}
