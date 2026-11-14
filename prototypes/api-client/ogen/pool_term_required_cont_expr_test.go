package ogenclient

import (
	"encoding/json"
	"testing"
)

// These cases exercise the real Pool.Term variant models generated from
// openapi/openapi.yaml, rather than the isolated recursion-control fixtures.
func TestPoolTermContExprIsOptional(t *testing.T) {
	tests := []struct {
		name       string
		newValue   func() any
		missing    string
		withCont   string
	}{
		{
			name:     "hire",
			newValue: func() any { return &PoolTermHireSpec{} },
			missing:  `{"conn_ph":"c","proc_qn":"p"}`,
			withCont: `{"conn_ph":"c","proc_qn":"p","cont_expr":{"kind":"fire","conn_ph":"c"}}`,
		},
		{
			name:     "apply",
			newValue: func() any { return &PoolTermApplySpec{} },
			missing:  `{"conn_ph":"c","proc_qn":"p"}`,
			withCont: `{"conn_ph":"c","proc_qn":"p","cont_expr":{"kind":"fire","conn_ph":"c"}}`,
		},
		{
			name:     "acquire",
			newValue: func() any { return &PoolTermAcquireSpec{} },
			missing:  `{"conn_ph":"c"}`,
			withCont: `{"conn_ph":"c","cont_expr":{"kind":"fire","conn_ph":"c"}}`,
		},
		{
			name:     "accept",
			newValue: func() any { return &PoolTermAcceptSpec{} },
			missing:  `{"conn_ph":"c"}`,
			withCont: `{"conn_ph":"c","cont_expr":{"kind":"fire","conn_ph":"c"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/missing cont_expr is accepted", func(t *testing.T) {
			if err := json.Unmarshal([]byte(tt.missing), tt.newValue()); err != nil {
				t.Fatalf("expected omitted optional cont_expr to decode: %v", err)
			}
		})

		t.Run(tt.name+"/present cont_expr is decoded", func(t *testing.T) {
			got := tt.newValue()
			if err := json.Unmarshal([]byte(tt.withCont), got); err != nil {
				t.Fatalf("expected present optional cont_expr to decode: %v", err)
			}
		})
	}
}
