package commexch

import (
	"github.com/huandu/go-sqlbuilder"

	"orglang/go-engine/adt/semcomm"
	"orglang/go-engine/pool/commturn"
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
	semBuilder := sqlbuilder.NewStruct(new(semcomm.SemRefDS)).For(sqlbuilder.PostgreSQL)
	connBuilder := sqlbuilder.NewStruct(new(exchRecDS)).For(sqlbuilder.PostgreSQL)
	stepBuilder := sqlbuilder.NewStruct(new(commturn.TurnRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{semBuilder, connBuilder, stepBuilder}
}

func (qb *sqlBuilder) insertRec(rec exchRecDS) (string, []any) {
	return qb.connBuilder.InsertInto(poolConns, rec).Build()
}

func (qb *sqlBuilder) updateRec(mod exchModDS) (string, []any) {
	conn := qb.connBuilder.Update(poolConns, mod)
	if mod.OffsetNr.Valid {
		conn.Set(conn.Assign("comm_on", mod.OffsetNr.V))
	}
	conn.Where(conn.Equal("comm_id", mod.CommID))
	return conn.Build()
}

func (qb *sqlBuilder) selectSnap(qry exchQryDS) (string, []any) {
	sem := qb.semBuilder.SelectFrom(commSems + "sem")
	step := qb.stepBuilder.SelectFrom(poolSteps + "step")
	step.Where(step.Equal("comm_id", qry.CommID))
	if qry.ChnlID.Valid {
		step.Where(step.Equal("chnl_id", qry.ChnlID.V))
	}
	steps := sqlbuilder.PostgreSQL.NewCTEBuilder()
	return sem.With(steps.With(sqlbuilder.CTEQuery("steps").As(step))).
		SelectMore(
			sem.BuilderAs(sqlbuilder.Build("SELECT array_agg(row(step.*)) FROM steps step WHERE step.comm_rn > conn.comm_on"), "steps"),
		).
		Join(poolConns+"conn", "conn.comm_id = sem.comm_id").
		Where(sem.Equal("sem.comm_id", qry.CommID)).
		Build()
}
