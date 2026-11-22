package adt

import "testing"

func TestNewAggregateRefs(t *testing.T) {
	t.Parallel()

	desc := NewDescRef()
	if desc.DescID.IsEmpty() || desc.DescRN == 0 {
		t.Fatalf("NewDescRef() = %#v, want non-empty identity and revision", desc)
	}

	comp := NewCompRef()
	if comp.CompID.IsEmpty() || comp.CompRN == 0 {
		t.Fatalf("NewCompRef() = %#v, want non-empty identity and revision", comp)
	}
}
