package comproot

import "testing"

func TestNew(t *testing.T) {
	t.Parallel()
	ref := New()
	if ref.CompID.IsEmpty() || ref.CompRN == 0 {
		t.Fatalf("New() = %#v, want non-empty identity and revision", ref)
	}
}
