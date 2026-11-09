package adapter

import (
	"github.com/huandu/go-sqlbuilder"

	commconn "orglang/go-engine/internal/pool/commconn/core"
	commturn "orglang/go-engine/internal/pool/commturn/core"
)

type sqlBuilder struct {
	connBuilder *sqlbuilder.Struct
	turnBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() commconn.QueryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	connBuilder := sqlbuilder.NewStruct(new(commconn.ConnRecDS)).For(sqlbuilder.PostgreSQL)
	turnBuilder := sqlbuilder.NewStruct(new(commturn.ActRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{connBuilder, turnBuilder}
}

func (qb *sqlBuilder) InsertRec(rec commconn.ConnRecDS) (string, []any) {
	return qb.connBuilder.InsertInto(compConns, rec).Build()
}

func (qb *sqlBuilder) InsertTurnRec(rec commturn.ActRecDS) (string, []any) {
	return qb.turnBuilder.InsertInto(connActs, rec).Build()
}

func (qb *sqlBuilder) UpdateRec(mod commconn.ConnModDS) (string, []any) {
	conn := qb.connBuilder.Update(compConns, mod)
	assignments := make([]string, 0, 2)
	if mod.OffsetNr.Valid {
		assignments = append(assignments, conn.Assign("offset_nr", mod.OffsetNr.V))
	}
	if mod.ActNr.Valid {
		assignments = append(assignments, conn.Assign("act_nr", mod.ActNr.V))
	}
	if len(assignments) > 0 {
		conn.Set(assignments...)
	}
	conn.Where(conn.Equal("conn_id", mod.ConnID))
	return conn.Build()
}

func (qb *sqlBuilder) SelectSnap(qry commconn.ConnQryDS) (string, []any) {
	conn := sqlbuilder.PostgreSQL.NewSelectBuilder().Select("conn.conn_id", "conn.provider_comp_id", "exec.comp_rn AS provider_comp_rn", "conn.offset_nr", "conn.act_nr").From(compConns+"conn").Join("pool_comp_execs exec", "exec.comp_id = conn.provider_comp_id")
	turn := qb.turnBuilder.SelectFrom(connActs + "turn")
	turn.Where(turn.Equal("conn_id", qry.ConnID))
	if qry.ChnlID.Valid {
		turn.Where(turn.Equal("chnl_id", qry.ChnlID.V))
	}
	turns := sqlbuilder.PostgreSQL.NewCTEBuilder()
	return conn.With(turns.With(sqlbuilder.CTEQuery("turns").As(turn))).
		SelectMore(
			conn.BuilderAs(sqlbuilder.Build("SELECT array_agg(row(t.*) ORDER BY t.act_nr) FROM turns t WHERE t.act_nr > conn.offset_nr"), "turns"),
		).
		Where(conn.Equal("conn.conn_id", qry.ConnID)).
		Build()
}
