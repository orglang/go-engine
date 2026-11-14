package ogenclient

import (
	"encoding/json"
	"testing"
)

// The required experiment preserves the wire-level required contract while
// pointer-boxing only changes the generated Go representation.
func TestPoolTermContExprIsRequired(t *testing.T) {
	tests := []struct {
		name     string
		newValue func() any
		missing  string
		present  string
	}{
		{"hire", func() any { return &PoolTermHireSpec{} }, `{"conn_ph":"c","proc_qn":"p"}`, `{"conn_ph":"c","proc_qn":"p","cont_expr":{"kind":"fire","conn_ph":"c"}}`},
		{"apply", func() any { return &PoolTermApplySpec{} }, `{"conn_ph":"c","proc_qn":"p"}`, `{"conn_ph":"c","proc_qn":"p","cont_expr":{"kind":"fire","conn_ph":"c"}}`},
		{"acquire", func() any { return &PoolTermAcquireSpec{} }, `{"conn_ph":"c"}`, `{"conn_ph":"c","cont_expr":{"kind":"fire","conn_ph":"c"}}`},
		{"accept", func() any { return &PoolTermAcceptSpec{} }, `{"conn_ph":"c"}`, `{"conn_ph":"c","cont_expr":{"kind":"fire","conn_ph":"c"}}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("missing is rejected", func(t *testing.T) {
				if err := json.Unmarshal([]byte(tt.missing), tt.newValue()); err == nil {
					t.Fatal("expected missing required cont_expr to be rejected")
				}
			})
			t.Run("present value round-trips", func(t *testing.T) {
				got := tt.newValue()
				if err := json.Unmarshal([]byte(tt.present), got); err != nil {
					t.Fatalf("decode present cont_expr: %v", err)
				}
				encoded, err := json.Marshal(got)
				if err != nil { t.Fatalf("encode present cont_expr: %v", err) }
				decoded := tt.newValue()
				if err := json.Unmarshal(encoded, decoded); err != nil {
					t.Fatalf("decode encoded value: %v", err)
				}
			})
			t.Run("explicit null is rejected", func(t *testing.T) {
				var payload map[string]json.RawMessage
				if err := json.Unmarshal([]byte(tt.missing), &payload); err != nil { t.Fatal(err) }
				payload["cont_expr"] = json.RawMessage("null")
				withNull, err := json.Marshal(payload)
				if err != nil { t.Fatal(err) }
				if err := json.Unmarshal(withNull, tt.newValue()); err == nil {
					t.Fatal("expected explicit null cont_expr to be rejected")
				}
			})
		})
	}
}

func TestRequiredPoolTermContExprNestedRecursion(t *testing.T) {
	payload := `{"conn_ph":"c","proc_qn":"p","cont_expr":{"kind":"hire","conn_ph":"c","proc_qn":"nested","cont_expr":{"kind":"fire","conn_ph":"c"}}}`
	got := &PoolTermHireSpec{}
	if err := json.Unmarshal([]byte(payload), got); err != nil {
		t.Fatalf("decode nested required cont_expr: %v", err)
	}
	encoded, err := json.Marshal(got)
	if err != nil { t.Fatalf("encode nested recursive value: %v", err) }
	decoded := &PoolTermHireSpec{}
	if err := json.Unmarshal(encoded, decoded); err != nil {
		t.Fatalf("decode nested round-trip: %v", err)
	}
}
