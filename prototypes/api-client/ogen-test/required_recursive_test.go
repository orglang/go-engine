package ogentest

import (
	"encoding/json"
	"testing"

	ogenclient "orglang/go-engine/prototypes/api-client/ogen"
)

func TestPoolTypeDefSnapExprRequired(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantErr bool
	}{
		{
			name:    "present",
			payload: `{"ref":{"type_id":"t","type_rn":1},"expr":{"kind":"one"}}`,
		},
		{
			name:    "missing",
			payload: `{"ref":{"type_id":"t","type_rn":1}}`,
			wantErr: true,
		},
		{
			name:    "null",
			payload: `{"ref":{"type_id":"t","type_rn":1},"expr":null}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got ogenclient.PoolTypeDefSnap
			err := json.Unmarshal([]byte(tt.payload), &got)
			if (err != nil) != tt.wantErr {
				t.Fatalf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && got.Expr == nil {
				t.Fatal("decoded required expr as nil")
			}
		})
	}
}

func TestPoolTypeLolliRequiredRecursiveFields(t *testing.T) {
	tests := []struct {
		name    string
		payload string
		wantErr bool
	}{
		{
			name:    "both recursive fields present",
			payload: `{"kind":"lolli","val_expr":{"kind":"one"},"cont_expr":{"kind":"lolli","val_expr":{"kind":"one"},"cont_expr":{"kind":"one"}}}`,
		},
		{
			name:    "missing val_expr",
			payload: `{"kind":"lolli","cont_expr":{"kind":"one"}}`,
			wantErr: true,
		},
		{
			name:    "null val_expr",
			payload: `{"kind":"lolli","val_expr":null,"cont_expr":{"kind":"one"}}`,
			wantErr: true,
		},
		{
			name:    "missing cont_expr",
			payload: `{"kind":"lolli","val_expr":{"kind":"one"}}`,
			wantErr: true,
		},
		{
			name:    "null cont_expr",
			payload: `{"kind":"lolli","val_expr":{"kind":"one"},"cont_expr":null}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got ogenclient.PoolTypeExprSpec
			err := json.Unmarshal([]byte(tt.payload), &got)
			if (err != nil) != tt.wantErr {
				t.Fatalf("json.Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
