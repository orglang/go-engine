package uniqsym

import (
	"database/sql"
	"strings"
	"testing"
)

var sunnyTests = []struct {
	name string
	str  string
	adt  ADT
}{
	{"sym only", "a", ADT{"a", nil}},
	{"sym and single-segment ns", "b.a", ADT{"a", &ADT{"b", nil}}},
	{"sym and two-segment ns", "c.b.a", ADT{"a", &ADT{"b", &ADT{"c", nil}}}},
}

func TestConvertFromStringSuccess(t *testing.T) {
	for _, test := range sunnyTests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ConvertFromString(test.str)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !got.Equal(test.adt) {
				t.Errorf("got %+v, want %+v", got, test.adt)
			}
		})
	}
}

func TestConvertToStringSuccess(t *testing.T) {
	for _, test := range sunnyTests {
		t.Run(test.name, func(t *testing.T) {
			got := ConvertToString(test.adt)
			if got != test.str {
				t.Errorf("got %s, want %s", got, test.str)
			}
		})
	}
}

func TestConvertFromStringError(t *testing.T) {
	rainyTests := []struct {
		name string
		str  string
	}{
		{"empty string", ""},
		{"separator only", "."},
		{"invalid symbol", "a."},
		{"invalid namespace", ".a"},
		{"whitespace", "a b"},
		{"slash", "a/b"},
		{"non-ascii", "привет"},
		{"empty segment", "a..b"},
		{"over maximum length", strings.Repeat("a", maxLength+1)},
	}
	for _, test := range rainyTests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := ConvertFromString(test.str); err == nil {
				t.Error("expected validation error")
			}
		})
	}
}

func TestConvertFromStringAcceptsMaximumLength(t *testing.T) {
	if _, err := ConvertFromString(strings.Repeat("a", maxLength)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConvertFromNullString(t *testing.T) {
	if got, err := ConvertFromNullString(sql.NullString{}); err != nil || got != empty {
		t.Fatalf("got (%v, %v), want (empty, nil)", got, err)
	}
	if _, err := ConvertFromNullString(sql.NullString{String: "bad value", Valid: true}); err == nil {
		t.Fatal("expected validation error")
	}
	if _, err := ConvertFromNullString(sql.NullString{String: "a.b", Valid: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
