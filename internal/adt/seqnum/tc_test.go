package seqnum

import "testing"

func TestConvertToInt(t *testing.T) {
	if got := ConvertToInt(ADT(7)); got != 7 {
		t.Fatalf("ConvertToInt() = %d, want 7", got)
	}
}
