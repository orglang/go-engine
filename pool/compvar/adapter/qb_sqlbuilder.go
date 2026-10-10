package adapter

import (
	"github.com/huandu/go-sqlbuilder"

	"orglang/go-engine/adt/compvar"

	poolcompvar "orglang/go-engine/pool/compvar/core"
)

type sqlBuilder struct {
	varBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() poolcompvar.QueryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	varBuilder := sqlbuilder.NewStruct(new(compvar.VarRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{varBuilder}
}

func (qb *sqlBuilder) InsertRec(table string, rec compvar.VarRecDS) (string, []any) {
	return qb.varBuilder.InsertInto(table, rec).Build()
}
