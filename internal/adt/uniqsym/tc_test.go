package uniqsym

import (
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

func TestConvertFromStringRejectsEmptySegments(t *testing.T) {
	for _, value := range []string{"", ".", "a.", ".a", "a..b"} {
		t.Run(value, func(t *testing.T) {
			if _, err := ConvertFromString(value); err == nil {
				t.Fatal("expected conversion error")
			}
		})
	}
}

func TestConvertFromStringDoesNotValidateQualifiedSymbolSyntax(t *testing.T) {
	for _, value := range []string{"a b", "a/b", "привет", strings.Repeat("a", 513)} {
		t.Run(value, func(t *testing.T) {
			if _, err := ConvertFromString(value); err != nil {
				t.Fatalf("unexpected constructor validation: %v", err)
			}
		})
	}
}

