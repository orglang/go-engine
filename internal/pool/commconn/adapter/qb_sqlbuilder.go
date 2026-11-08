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
	turnBuilder := sqlbuilder.NewStruct(new(commturn.TurnRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{connBuilder, turnBuilder}
}

func (qb *sqlBuilder) InsertRec(rec commconn.ConnRecDS) (string, []any) {
	return qb.connBuilder.InsertInto(commConns, rec).Build()
}

func (qb *sqlBuilder) InsertTurnRec(rec commturn.TurnRecDS) (string, []any) {
	return qb.turnBuilder.InsertInto(commTurns, rec).Build()
}

func (qb *sqlBuilder) UpdateRec(mod commconn.ConnModDS) (string, []any) {
	conn := qb.connBuilder.Update(commConns, mod)
	if mod.OffsetNr.Valid {
		conn.Set(conn.Assign("offset_nr", mod.OffsetNr.V))
	}
	conn.Where(conn.Equal("comm_id", mod.ConnID))
	return conn.Build()
}

func (qb *sqlBuilder) SelectSnap(qry commconn.ConnQryDS) (string, []any) {
	conn := qb.connBuilder.SelectFrom(commConns + "conn")
	turn := qb.turnBuilder.SelectFrom(commTurns + "turn")
	turn.Where(turn.Equal("comm_id", qry.ConnID))
	if qry.ChnlID.Valid {
		turn.Where(turn.Equal("chnl_id", qry.ChnlID.V))
	}
	turns := sqlbuilder.PostgreSQL.NewCTEBuilder()
	return conn.With(turns.With(sqlbuilder.CTEQuery("turns").As(turn))).
		SelectMore(
			conn.BuilderAs(sqlbuilder.Build("SELECT array_agg(row(t.*)) FROM turns t WHERE t.comm_rn > conn.offset_nr"), "turns"),
		).
		Where(conn.Equal("conn.comm_id", qry.ConnID)).
		Build()
}
