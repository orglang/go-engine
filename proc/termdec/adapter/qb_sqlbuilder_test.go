package adapter

import (
	"reflect"
	"testing"

	"orglang/go-engine/adt/termvar"

	termdec "orglang/go-engine/proc/termdec/core"
)

func TestInsertRec(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.InsertRec(termdec.DecRecDS{})
	wantSQL := "INSERT INTO proc_term_decs  (term_id, term_rn, liab_var, asset_vars) VALUES ($1, $2, $3, $4)"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{"", int64(0), termvar.VarRecDS{}, []termvar.VarRecDS(nil)}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}
