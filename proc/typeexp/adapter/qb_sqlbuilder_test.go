package adapter

import (
	"reflect"
	"testing"

	proctypeexp "orglang/go-engine/proc/typeexp/core"
)

func TestInsertRec(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.InsertRec(proctypeexp.StateDS{})
	wantSQL := "INSERT INTO pool_type_exps (exp_vk, sup_exp_vk, kind, spec) " +
		"VALUES ($1, $2, $3, $4) ON CONFLICT (exp_vk) DO NOTHING"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{int64(0), int64(0), proctypeexp.ExpKindDS(0), proctypeexp.ExpSpecDS{}}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}
