package contkey

import (
	"bytes"
	"testing"
)

func TestHashJSONCanonicalizesObjectMemberOrder(t *testing.T) {
	left := map[string]any{"z": int64(1), "a": "value"}
	right := map[string]any{"a": "value", "z": int64(1)}

	leftKey, leftJSON, err := HashJSON(left)
	if err != nil {
		t.Fatal(err)
	}
	rightKey, rightJSON, err := HashJSON(right)
	if err != nil {
		t.Fatal(err)
	}

	if leftKey != rightKey {
		t.Fatalf("keys differ: %d != %d", leftKey, rightKey)
	}
	if !bytes.Equal(leftJSON, rightJSON) {
		t.Fatalf("canonical JSON differs: %s != %s", leftJSON, rightJSON)
	}
}

func TestHashJSONPreservesInt64Precision(t *testing.T) {
	const maxInt64 = int64(1<<63 - 1)
	key, canonical, err := HashJSON(map[string]any{"value": maxInt64})
	if err != nil {
		t.Fatal(err)
	}
	if key == Zero {
		t.Fatal("expected non-zero content key")
	}
	const want = `{"value":9223372036854775807}`
	if string(canonical) != want {
		t.Fatalf("canonical JSON = %s, want %s", canonical, want)
	}
}
