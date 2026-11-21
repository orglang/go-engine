package adapter

import (
	"database/sql"
	"testing"

	"orglang/go-engine/internal/adt/uniqsym"
)

func TestConvertFromNullString(t *testing.T) {
	var empty uniqsym.ADT
	if got, err := ConvertFromNullString(sql.NullString{}); err != nil || got != empty {
		t.Fatalf("got (%v, %v), want empty, nil", got, err)
	}
	if _, err := ConvertFromNullString(sql.NullString{String: "bad value", Valid: true}); err != nil {
		t.Fatalf("unexpected constructor validation: %v", err)
	}
	if got, err := ConvertFromNullString(sql.NullString{String: "a.b", Valid: true}); err != nil || !got.Equal(mustUniqsym(t, "a.b")) {
		t.Fatalf("got (%v, %v), want a.b, nil", got, err)
	}
}

func mustUniqsym(t *testing.T, value string) uniqsym.ADT {
	t.Helper()
	got, err := uniqsym.ConvertFromString(value)
	if err != nil {
		t.Fatal(err)
	}
	return got
}
