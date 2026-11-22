package descroot

import "testing"

func TestNew(t *testing.T) {
	t.Parallel()
	ref := New()
	if ref.DescID.IsEmpty() || ref.DescRN == 0 {
		t.Fatalf("New() = %#v, want non-empty identity and revision", ref)
	}
}
