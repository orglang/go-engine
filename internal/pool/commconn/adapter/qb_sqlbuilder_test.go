package adapter

import (
	"database/sql"
	"slices"
	"testing"

	commconn "orglang/go-engine/internal/pool/commconn/core"
	commturn "orglang/go-engine/internal/pool/commturn/core"
)

func TestInsertRec(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.InsertRec(commconn.ConnRecDS{})
	wantSQL := "INSERT INTO pool_connections  (conn_id, conn_rn, offset_nr) VALUES ($1, $2, $3)"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{"", int64(0), int64(0)}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}

func TestInsertTurnRec(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.InsertTurnRec(commturn.TurnRecDS{})
	wantSQL := "INSERT INTO pool_connection_acts  (conn_id, conn_rn, comp_id, chnl_id, kind, exp) VALUES ($1, $2, $3, $4, $5, $6)"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	if len(gotArgs) != 6 {
		t.Errorf("got %d args, want 6: %v", len(gotArgs), gotArgs)
	}
}

func TestUpdateRec(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.UpdateRec(commconn.ConnModDS{OffsetNr: sql.Null[int64]{Valid: true}})
	wantSQL := "UPDATE pool_connections  SET offset_nr = $1 WHERE conn_id = $2 AND conn_rn = $3 AND offset_nr = $4"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{int64(0), "", int64(0), int64(0)}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}

func TestSelectSnap(t *testing.T) {
	t.Parallel()
	qb := newSQLBuilder()
	gotSQL, gotArgs := qb.SelectSnap(commconn.ConnQryDS{ChnlID: sql.Null[string]{Valid: true}})
	wantSQL := "WITH turns AS (SELECT turn.conn_id, turn.conn_rn, turn.comp_id, turn.chnl_id, turn.kind, turn.exp " +
		"FROM pool_connection_acts turn WHERE conn_id = $1 AND chnl_id = $2) " +
		"SELECT conn.conn_id, conn.conn_rn, conn.offset_nr, " +
		"(SELECT array_agg(row(t.*) ORDER BY t.conn_rn) FROM turns t WHERE t.conn_rn > conn.offset_nr) AS turns " +
		"FROM pool_connections conn WHERE conn.conn_id = $3"
	if gotSQL != wantSQL {
		t.Errorf("got sql %q, want %q", gotSQL, wantSQL)
	}
	wantArgs := []any{"", "", ""}
	if !slices.Equal(gotArgs, wantArgs) {
		t.Errorf("got args %v, want %v", gotArgs, wantArgs)
	}
}
