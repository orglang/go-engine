package adapter

import (
	"strings"
	"testing"
)

func TestDecSpecVPValidate(t *testing.T) {
	tests := []struct {
		name  string
		value string
		valid bool
}{
		{name: "simple", value: "alpha", valid: true},
		{name: "qualified", value: "alpha.beta_2", valid: true},
		{name: "maximum length", value: strings.Repeat("a", 512), valid: true},
		{name: "empty", value: "", valid: false},
		{name: "empty segment", value: "alpha..beta", valid: false},
		{name: "whitespace", value: "alpha beta", valid: false},
		{name: "slash", value: "alpha/beta", valid: false},
		{name: "too long", value: strings.Repeat("a", 513), valid: false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dto := DecSpecVP{TermQN: test.value}
			err := dto.Validate()
			if test.valid && err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
			if !test.valid && err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}
