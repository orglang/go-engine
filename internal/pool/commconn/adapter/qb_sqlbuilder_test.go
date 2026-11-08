package adapter

import (
	"database/sql"
	"slices"
	"testing"

	commconn "orglang/go-engine/internal/pool/commconn/core"
)

func TestInsertRec(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.InsertRec(commconn.ConnRecDS{})
	wantSQL := "INSERT INTO pool_comm_conns  (comm_id, comm_rn, offset_nr) VALUES ($1, $2, $3)"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{"", int64(0), int64(0)}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}

func TestUpdateRec(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.UpdateRec(commconn.ConnModDS{OffsetNr: sql.Null[int64]{Valid: true}})
	wantSQL := "UPDATE pool_comm_conns  SET offset_nr = $1 WHERE comm_id = $2"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{int64(0), ""}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}

func TestSelectSnap(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.SelectSnap(commconn.ConnQryDS{ChnlID: sql.Null[string]{Valid: true}})
	wantSQL := "WITH turns AS (SELECT turn.comm_id, turn.comm_rn, turn.comp_id, turn.chnl_id, turn.kind, turn.exp " +
		"FROM pool_comm_turns turn WHERE comm_id = $1 AND chnl_id = $2) " +
		"SELECT conn.comm_id, conn.comm_rn, conn.offset_nr, " +
		"(SELECT array_agg(row(t.*)) FROM turns t WHERE t.comm_rn > conn.offset_nr) AS turns " +
		"FROM pool_comm_conns conn WHERE conn.comm_id = $3"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{"", "", ""}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}
