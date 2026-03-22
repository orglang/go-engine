package commsem

import (
	"github.com/huandu/go-sqlbuilder"
)

type sqlBuilder struct {
	semBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuikder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	semBuilder := sqlbuilder.NewStruct(new(semRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{semBuilder}
}

func (qb *sqlBuilder) insertRec(rec semRecDS) (string, []any) {
	return qb.semBuilder.InsertInto("comm_sems", rec).Build()
}
