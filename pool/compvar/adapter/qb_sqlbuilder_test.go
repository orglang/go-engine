package adapter

import (
	"database/sql"
	"slices"
	"testing"

	"orglang/go-engine/adt/compvar"
)

func TestInsert(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	rec := compvar.VarRecDS{
		CompID: sql.NullString{String: "comp-id-1", Valid: true},
		CompRN: sql.NullInt64{Int64: 7, Valid: true},
		CommID: sql.NullString{String: "comm-id-1", Valid: true},
		ChnlID: sql.NullString{String: "chnl-id-1", Valid: true},
		ChnlPH: sql.NullString{String: "chnl-1", Valid: true},
		ExpVK:  sql.NullInt64{Int64: 1, Valid: true},
		ChnlBS: sql.NullInt16{Int16: 1, Valid: true},
	}
	gotSQL, gotArgs := qb.InsertRec(poolStructVars, rec)
	wantSQL := "INSERT INTO pool_struct_vars (comp_id, comp_rn, comm_id, chnl_id, chnl_ph, exp_vk, side) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7)"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{rec.CompID, rec.CompRN, rec.CommID, rec.ChnlID, rec.ChnlPH, rec.ExpVK, rec.ChnlBS}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}
