package poolvar

import (
	"orglang/go-engine/adt/implvar"

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
	varBuilder := sqlbuilder.NewStruct(new(implvar.VarRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{varBuilder}
}

func (qb *sqlBuilder) insertRec(table string, rec implvar.VarRecDS) (string, []any) {
	return qb.semBuilder.InsertInto(table, rec).Build()
}
