package arch

import (
	"strings"
	"testing"
)

func TestCheck(t *testing.T) {
	imports := map[string][]string{
		"orglang/go-engine/pool/typedef/core": {
			"orglang/go-engine/pool/typeexp/core",
			"orglang/go-engine/adt/uniqsym",
			"orglang/go-engine/lib/db",
			"orglang/go-engine/pool/typedef/adapter",
		},
		"orglang/go-engine/pool/typedef/adapter": {
			"orglang/go-engine/pool/typedef/core",
		},
		"orglang/go-engine/proc/compexec/core": {
			"orglang/go-engine/proc/termdec/core",
		},
	}
	want := []Violation{{
		From: "orglang/go-engine/pool/typedef/core",
		To:   "orglang/go-engine/pool/typedef/adapter",
	}}
	got := Check(imports)
	if len(got) != len(want) {
		t.Fatalf("got %d violations, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("violation %d = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestCheckIgnoresNonAgnosticSources(t *testing.T) {
	imports := map[string][]string{
		"orglang/go-engine/pool/typedef/adapter": {
			"orglang/go-engine/lib/db/ds_pgx",
		},
		"orglang/go-engine/lib/ws": {
			"orglang/go-engine/pool/compexec/adapter",
		},
		"orglang/go-engine/adt/valkey": {
			"orglang/go-engine/adt/compsem/adapter",
		},
	}
	if got := Check(imports); len(got) != 0 {
		t.Errorf("got %v, want no violations", got)
	}
}

func TestCheckMessageNamesBothEdges(t *testing.T) {
	v := Violation{
		From: "orglang/go-engine/pool/typedef/core",
		To:   "orglang/go-engine/pool/typedef/adapter",
	}
	msg := v.String()
	for _, want := range []string{"pool/typedef/core", "pool/typedef/adapter"} {
		if !strings.Contains(msg, want) {
			t.Errorf("message %q missing %q", msg, want)
		}
	}
}
