package progroot

import "testing"

func TestNew(t *testing.T) {
	t.Parallel()
	ref := New()
	if ref.ProgID.IsEmpty() || ref.ProgRN == 0 {
		t.Fatalf("New() = %#v, want non-empty identity and revision", ref)
	}
}
