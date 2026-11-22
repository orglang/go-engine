package reference

import "testing"

func TestNewIDsAreTypedAndNonEmpty(t *testing.T) {
	t.Parallel()

	if NewDescID().IsEmpty() {
		t.Fatal("NewDescID returned an empty ID")
	}
	if NewTypeID().IsEmpty() {
		t.Fatal("NewTypeID returned an empty ID")
	}
	if NewTermID().IsEmpty() {
		t.Fatal("NewTermID returned an empty ID")
	}
	if NewCompID().IsEmpty() {
		t.Fatal("NewCompID returned an empty ID")
	}
}

func TestAggregateRevisionsAdvanceIndependently(t *testing.T) {
	t.Parallel()

	descRN := NewDescRN()
	compRN := NewCompRN()
	if got := descRN.Next(); got != descRN+1 {
		t.Fatalf("DescRN.Next() = %v, want %v", got, descRN+1)
	}
	if got := compRN.Next(); got != compRN+1 {
		t.Fatalf("CompRN.Next() = %v, want %v", got, compRN+1)
	}
}
