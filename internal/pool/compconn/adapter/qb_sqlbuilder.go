package adapter

import (
	"github.com/huandu/go-sqlbuilder"

	compconn "orglang/go-engine/internal/pool/compconn/core"
	connact "orglang/go-engine/internal/pool/connact/core"
)

type sqlBuilder struct {
	connBuilder *sqlbuilder.Struct
	actBuilder  *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() compconn.QueryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	connBuilder := sqlbuilder.NewStruct(new(compconn.ConnRecDS)).For(sqlbuilder.PostgreSQL)
	actBuilder := sqlbuilder.NewStruct(new(connact.ActRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{connBuilder, actBuilder}
}

func (qb *sqlBuilder) InsertRec(rec compconn.ConnRecDS) (string, []any) {
	return qb.connBuilder.InsertInto(compConns, rec).Build()
}

func (qb *sqlBuilder) InsertActRec(rec connact.ActRecDS) (string, []any) {
	return qb.actBuilder.InsertInto(connActs, rec).Build()
}

func (qb *sqlBuilder) UpdateRec(mod compconn.ConnModDS) (string, []any) {
	conn := qb.connBuilder.Update(compConns, mod)
	assignments := make([]string, 0, 2)
	if mod.OffsetNr.Valid {
		assignments = append(assignments, compconn.Assign("offset_nr", mod.OffsetNr.V))
	}
	if mod.ActNr.Valid {
		assignments = append(assignments, compconn.Assign("act_nr", mod.ActNr.V))
	}
	if len(assignments) > 0 {
		compconn.Set(assignments...)
	}
	compconn.Where(compconn.Equal("conn_id", mod.ConnID))
	return compconn.Build()
}

func (qb *sqlBuilder) SelectSnap(qry compconn.ConnQryDS) (string, []any) {
	conn := sqlbuilder.PostgreSQL.NewSelectBuilder().Select("compconn.conn_id", "compconn.provider_comp_id", "exec.comp_rn AS provider_comp_rn", "compconn.offset_nr", "compconn.act_nr").From(compConns+"conn").Join("pool_comp_execs exec", "exec.comp_id = compconn.provider_comp_id")
	act := qb.actBuilder.SelectFrom(connActs + "act")
	act.Where(act.Equal("conn_id", qry.ConnID))
	if qry.ChnlID.Valid {
		act.Where(act.Equal("chnl_id", qry.ChnlID.V))
	}
	acts := sqlbuilder.PostgreSQL.NewCTEBuilder()
	return compconn.With(acts.With(sqlbuilder.CTEQuery("acts").As(act))).
		SelectMore(
			compconn.BuilderAs(sqlbuilder.Build("SELECT array_agg(row(t.*) ORDER BY t.act_nr) FROM acts t WHERE t.act_nr > compconn.offset_nr"), "acts"),
		).
		Where(compconn.Equal("compconn.conn_id", qry.ConnID)).
		Build()
}
