package termvar

import (
	"encoding/json"
	"testing"
)

func TestVarRecDSKeepsLegacyContentKeyJSONName(t *testing.T) {
	const payload = `{"id":"type-1","ph":"channel-1","vk":42}`

	var rec VarRecDS
	if err := json.Unmarshal([]byte(payload), &rec); err != nil {
		t.Fatalf("unmarshal VarRecDS: %v", err)
	}
	if rec.ExpCK != 42 {
		t.Fatalf("ExpCK = %d, want 42", rec.ExpCK)
	}

	got, err := json.Marshal(rec)
	if err != nil {
		t.Fatalf("marshal VarRecDS: %v", err)
	}
	if string(got) != payload {
		t.Fatalf("JSON = %s, want %s", got, payload)
	}
}
