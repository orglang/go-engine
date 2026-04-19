package commturn

import (
	"github.com/huandu/go-sqlbuilder"

	"orglang/go-engine/adt/termsem"
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
	semBuilder := sqlbuilder.NewStruct(new(termsem.SemRefDS)).For(sqlbuilder.PostgreSQL)
	stepBuilder := sqlbuilder.NewStruct(new(TurnRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{semBuilder, stepBuilder}
}

func (qb *sqlBuilder) insertRec(rec TurnRecDS) (string, []any) {
	return qb.stepBuilder.InsertInto(commTurns, rec).Build()
}
