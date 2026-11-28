package adapter

import (
	"database/sql"
	"testing"

	"orglang/go-engine/internal/adt/option"
	"orglang/go-engine/internal/adt/seqnum"
)

func TestConvertOptionToNullInt(t *testing.T) {
	tests := []struct {
		name string
		in   option.ADT[seqnum.ADT]
		want sql.Null[int64]
	}{
		{name: "nil"},
		{name: "none", in: option.None[seqnum.ADT]()},
		{name: "some", in: option.Some(seqnum.ADT(7)), want: sql.Null[int64]{V: 7, Valid: true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertOptionToNullInt(tt.in); got != tt.want {
				t.Fatalf("ConvertOptionToNullInt() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
