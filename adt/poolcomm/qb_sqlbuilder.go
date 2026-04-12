package poolcomm

import (
	"github.com/huandu/go-sqlbuilder"

	"orglang/go-engine/adt/commsem"
	"orglang/go-engine/adt/poolstep"
)

type sqlBuilder struct {
	commJoinBuilder *sqlbuilder.Struct
	semBuilder      *sqlbuilder.Struct
	connBuilder     *sqlbuilder.Struct
	stepBuilder     *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuikder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	commJoinBuilder := sqlbuilder.NewStruct(new(commJoinDS)).For(sqlbuilder.PostgreSQL)
	semBuilder := sqlbuilder.NewStruct(new(commsem.SemRefDS)).For(sqlbuilder.PostgreSQL)
	connBuilder := sqlbuilder.NewStruct(new(connRecDS)).For(sqlbuilder.PostgreSQL)
	stepBuilder := sqlbuilder.NewStruct(new(poolstep.StepRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{commJoinBuilder, semBuilder, connBuilder, stepBuilder}
}

func (qb *sqlBuilder) insertRec(rec connRecDS) (string, []any) {
	return qb.connBuilder.InsertInto(poolConns, rec).Build()
}

func (qb *sqlBuilder) updateRec(mod commModDS) (string, []any) {
	conn := qb.connBuilder.Update(poolConns, mod)
	if mod.CommON.Valid {
		conn.Set(conn.Assign("comm_on", mod.CommON.V))
	}
	conn.Where(conn.Equal("comm_id", mod.CommID))
	return conn.Build()
}

func (qb *sqlBuilder) selectSnap(qry commQryDS) (string, []any) {
	comm := qb.commJoinBuilder.SelectFrom(commSems + "sem")
	step := qb.stepBuilder.SelectFrom(poolSteps + "s")
	step.Where(step.Equal("comm_id", qry.CommID))
	if qry.ChnlID.Valid {
		step.Where(step.Equal("chnl_id", qry.ChnlID.V))
	}
	steps := sqlbuilder.PostgreSQL.NewCTEBuilder()
	return comm.With(steps.With(sqlbuilder.CTEQuery("steps").As(step))).
		SelectMore(
			comm.BuilderAs(sqlbuilder.Build("SELECT array_agg(row(step.*)) FROM steps step WHERE step.comm_rn > conn.comm_on"), "steps"),
		).
		Join(poolConns+"conn", "conn.comm_id = sem.comm_id").
		Where(comm.Equal("sem.comm_id", qry.CommID)).
		Build()
}
