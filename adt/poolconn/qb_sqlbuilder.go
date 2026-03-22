package poolconn

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
	semBuilder := sqlbuilder.NewStruct(new(connRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{semBuilder}
}

func (qb *sqlBuilder) insertRec(rec connRecDS) (string, []any) {
	return qb.semBuilder.InsertInto("pool_conns", rec).Build()
}
