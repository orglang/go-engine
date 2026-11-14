package probe

import (
	"encoding/json"
	"testing"
)

func TestRequiredNullableField(t *testing.T) {
	t.Run("missing field is rejected", func(t *testing.T) {
		var got B
		if err := json.Unmarshal([]byte(`{}`), &got); err == nil {
			t.Fatal("expected missing required field a to be rejected")
		}
	})

	t.Run("explicit null is accepted", func(t *testing.T) {
		var got B
		if err := json.Unmarshal([]byte(`{"a":null}`), &got); err != nil {
			t.Fatalf("explicit null for required nullable field a must be accepted: %v", err)
		}
		if got.A != nil {
			t.Fatalf("expected A to remain nil for explicit null, got %#v", got.A)
		}
	})

	t.Run("non-null value is decoded", func(t *testing.T) {
		var got B
		if err := json.Unmarshal([]byte(`{"a":{"b":{"a":null}}}`), &got); err != nil {
			t.Fatalf("expected non-null recursive value to decode: %v", err)
		}
		if got.A == nil || got.A.B.A != nil {
			t.Fatalf("unexpected decoded recursive value: %#v", got)
		}
	})
}
