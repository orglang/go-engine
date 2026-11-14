package ogenclient

import (
	"encoding/json"
	"testing"
)

// This experiment deliberately weakens the Pool.Term wire contract.
// Missing cont_expr must therefore decode successfully in all four variants.
func TestOptionalPoolTermContExpr(t *testing.T) {
	tests := []struct {
		name string
		newValue func() any
		payload string
	}{
		{"hire", func() any { return &PoolTermHireSpec{} }, `{"conn_ph":"c","proc_qn":"p"}`},
		{"apply", func() any { return &PoolTermApplySpec{} }, `{"conn_ph":"c","proc_qn":"p"}`},
		{"acquire", func() any { return &PoolTermAcquireSpec{} }, `{"conn_ph":"c"}`},
		{"accept", func() any { return &PoolTermAcceptSpec{} }, `{"conn_ph":"c"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := json.Unmarshal([]byte(tt.payload), tt.newValue()); err != nil {
				t.Fatalf("expected missing optional cont_expr to decode: %v", err)
			}
		})
	}
}
