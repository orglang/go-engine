package execroot

import "testing"

func TestNew(t *testing.T) {
	t.Parallel()
	ref := New()
	if ref.ExecID.IsEmpty() || ref.ExecRN == 0 {
		t.Fatalf("New() = %#v, want non-empty identity and revision", ref)
	}
}
