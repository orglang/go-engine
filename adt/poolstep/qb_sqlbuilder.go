package poolstep

import (
	"github.com/huandu/go-sqlbuilder"

	"orglang/go-engine/adt/implsem"
)

type sqlBuilder struct {
	semBuilder  *sqlbuilder.Struct
	stepBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuikder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	semBuilder := sqlbuilder.NewStruct(new(implsem.SemRefDS)).For(sqlbuilder.PostgreSQL)
	stepBuilder := sqlbuilder.NewStruct(new(StepRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{semBuilder, stepBuilder}
}

func (qb *sqlBuilder) insertRec(rec StepRecDS) (string, []any) {
	return qb.stepBuilder.InsertInto(poolSteps, rec).Build()
}
