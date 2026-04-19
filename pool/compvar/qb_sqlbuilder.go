package compvar

import (
	"orglang/go-engine/adt/compvar"

	"github.com/huandu/go-sqlbuilder"
)

type sqlBuilder struct {
	varBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	varBuilder := sqlbuilder.NewStruct(new(compvar.VarRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{varBuilder}
}

func (qb *sqlBuilder) insertRec(table string, rec compvar.VarRecDS) (string, []any) {
	return qb.varBuilder.InsertInto(table, rec).Build()
}
