package adapter

import (
	"encoding/json"
	"testing"
)

func TestVarRecKeepsLegacyContentKeyJSONName(t *testing.T) {
	const payload = `{"id":"type-1","ph":"channel-1","vk":42}`

	var rec VarRec
	if err := json.Unmarshal([]byte(payload), &rec); err != nil {
		t.Fatalf("unmarshal VarRec: %v", err)
	}
	if rec.ExpCK != 42 {
		t.Fatalf("ExpCK = %d, want 42", rec.ExpCK)
	}

	got, err := json.Marshal(rec)
	if err != nil {
		t.Fatalf("marshal VarRec: %v", err)
	}
	if string(got) != payload {
		t.Fatalf("JSON = %s, want %s", got, payload)
	}
}
