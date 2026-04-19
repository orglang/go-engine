package commexch

import (
	"github.com/huandu/go-sqlbuilder"

	"orglang/go-engine/pool/commturn"
)

type sqlBuilder struct {
	exchBuilder *sqlbuilder.Struct
	turnBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	exchBuilder := sqlbuilder.NewStruct(new(exchRecDS)).For(sqlbuilder.PostgreSQL)
	turnBuilder := sqlbuilder.NewStruct(new(commturn.TurnRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{exchBuilder, turnBuilder}
}

func (qb *sqlBuilder) insertRec(rec exchRecDS) (string, []any) {
	return qb.exchBuilder.InsertInto(commExchs, rec).Build()
}

func (qb *sqlBuilder) updateRec(mod exchModDS) (string, []any) {
	conn := qb.exchBuilder.Update(commExchs, mod)
	if mod.OffsetNr.Valid {
		conn.Set(conn.Assign("offset_nr", mod.OffsetNr.V))
	}
	conn.Where(conn.Equal("comm_id", mod.CommID))
	return conn.Build()
}

func (qb *sqlBuilder) selectSnap(qry exchQryDS) (string, []any) {
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
