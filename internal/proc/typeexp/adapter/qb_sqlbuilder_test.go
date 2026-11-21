package adapter

import (
	"reflect"
	"testing"
)

func TestInsertRec(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.InsertRec(State{})
	wantSQL := "INSERT INTO pool_type_exps (exp_ck, sup_exp_ck, kind, spec) " +
		"VALUES ($1, $2, $3, $4) ON CONFLICT (exp_ck) DO NOTHING"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{int64(0), int64(0), ExpKind(0), ExpSpec{}}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}
