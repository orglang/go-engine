package adapter

import (
	"reflect"
	"testing"

	termvaradapter "orglang/go-engine/internal/adt/termvar/adapter"

)

func TestInsertRec(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.InsertRec(DecRecDS{})
	wantSQL := "INSERT INTO proc_term_decs  (term_id, term_rn, liab_var, asset_vars) VALUES ($1, $2, $3, $4)"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{"", int64(0), termvaradapter.VarRecDS{}, []termvaradapter.VarRecDS(nil)}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}
