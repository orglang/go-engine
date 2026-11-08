package adapter

import (
	"github.com/huandu/go-sqlbuilder"

	commexch "orglang/go-engine/pool/commexch/core"
	commturn "orglang/go-engine/pool/commturn/core"
)

type sqlBuilder struct {
	exchBuilder *sqlbuilder.Struct
	turnBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() commexch.QueryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	exchBuilder := sqlbuilder.NewStruct(new(commexch.ExchRecDS)).For(sqlbuilder.PostgreSQL)
	turnBuilder := sqlbuilder.NewStruct(new(commturn.TurnRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{exchBuilder, turnBuilder}
}

func (qb *sqlBuilder) InsertRec(rec commexch.ExchRecDS) (string, []any) {
	return qb.exchBuilder.InsertInto(commExchs, rec).Build()
}

func (qb *sqlBuilder) UpdateRec(mod commexch.ExchModDS) (string, []any) {
	conn := qb.exchBuilder.Update(commExchs, mod)
	if mod.OffsetNr.Valid {
		conn.Set(conn.Assign("offset_nr", mod.OffsetNr.V))
	}
	conn.Where(conn.Equal("comm_id", mod.CommID))
	return conn.Build()
}

func (qb *sqlBuilder) SelectSnap(qry commexch.ExchQryDS) (string, []any) {
	exch := qb.exchBuilder.SelectFrom(commExchs + "exch")
	turn := qb.turnBuilder.SelectFrom(commTurns + "turn")
	turn.Where(turn.Equal("comm_id", qry.CommID))
	if qry.ChnlID.Valid {
		turn.Where(turn.Equal("chnl_id", qry.ChnlID.V))
	}
	turns := sqlbuilder.PostgreSQL.NewCTEBuilder()
	return exch.With(turns.With(sqlbuilder.CTEQuery("turns").As(turn))).
		SelectMore(
			exch.BuilderAs(sqlbuilder.Build("SELECT array_agg(row(t.*)) FROM turns t WHERE t.comm_rn > exch.offset_nr"), "turns"),
		).
		Where(exch.Equal("exch.comm_id", qry.CommID)).
		Build()
}
