package arch

import "testing"

func TestClassify(t *testing.T) {
	tests := []struct {
		path string
		want Kind
	}{
		{"orglang/go-engine/pool/typedef/core", Agnostic},
		{"orglang/go-engine/pool/typedef/domain", Agnostic},
		{"orglang/go-engine/proc/compexec/core", Agnostic},
		{"orglang/go-engine/pool/typedef/adapter", Specific},
		{"orglang/go-engine/pool/typedef/adapters/in/http", Specific},
		{"orglang/go-engine/pool/typedef/adapter/ds_pgx", Specific},
		{"orglang/go-engine/lib/db", Framework},
		{"orglang/go-engine/lib/ws", Framework},
		{"orglang/go-engine/adt/valkey", Shared},
		{"orglang/go-engine/adt/compsem", Shared},
		{"orglang/go-engine/pool/typedef", Unknown},
		{"orglang/go-engine/other/thing", Unknown},
	}
	for _, tt := range tests {
		if got := Classify(tt.path); got != tt.want {
			t.Errorf("Classify(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}
