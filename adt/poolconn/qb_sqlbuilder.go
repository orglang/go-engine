package poolconn

import (
	"github.com/huandu/go-sqlbuilder"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/poolstep"
)

type sqlBuilder struct {
	semBuilder  *sqlbuilder.Struct
	connBuilder *sqlbuilder.Struct
	stepBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuikder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	semBuilder := sqlbuilder.NewStruct(new(commsem.SemRefDS)).For(sqlbuilder.PostgreSQL)
	connBuilder := sqlbuilder.NewStruct(new(connRecDS)).For(sqlbuilder.PostgreSQL)
	stepBuilder := sqlbuilder.NewStruct(new(poolstep.StepRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{semBuilder, connBuilder, stepBuilder}
}

func (qb *sqlBuilder) insertRec(rec connRecDS) (string, []any) {
	return qb.connBuilder.InsertInto(poolConns, rec).Build()
}

func (qb *sqlBuilder) updateRec(mod connModDS) (string, []any) {
	conn := qb.connBuilder.Update(poolConns, mod)
	if mod.CommON.Valid {
		conn.Set(conn.Assign("comm_on", mod.CommON.V))
	}
	conn.Where(conn.Equal("comm_id", mod.CommID))
	return conn.Build()
}

func (qb *sqlBuilder) selectSnap(qry connQryDS) (string, []any) {
	step := sqlbuilder.PostgreSQL.NewSelectBuilder()
	step.SQL("SELECT array_agg(row(step.*))")
	step.From(poolSteps + " step")
	step.Where(step.Equal("step.comm_id", qry.CommID))
	if qry.ChnlID.Valid {
		step.Where(step.Equal("step.chnl_id", qry.ChnlID.V))
	}
	sem := qb.semBuilder.SelectFrom(commSems + " sem")
	sem.SelectMore(sem.BuilderAs(step, "steps"))
	sem.Where(sem.Equal("sem.comm_id", qry.CommID))
	return sem.Build()
}
