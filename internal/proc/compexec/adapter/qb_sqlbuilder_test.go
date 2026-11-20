package adapter

import (
	"reflect"
	"testing"

)

func TestInsertRec(t *testing.T) {
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.InsertRec(ExecRecDS{})
	wantSQL := "INSERT INTO proc_comp_execs  (comp_id, comp_rn, liab_mode) VALUES ($1, $2, $3)"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{"", int64(0), int16(0)}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}
