package adapter

import (
	"github.com/huandu/go-sqlbuilder"

	adtcompvar "orglang/go-engine/adt/compvar"

	compvar "orglang/go-engine/pool/compvar/core"
)

type sqlBuilder struct {
	varBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() compvar.QueryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	varBuilder := sqlbuilder.NewStruct(new(adtcompvar.VarRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{varBuilder}
}

func (qb *sqlBuilder) InsertRec(table string, rec adtcompvar.VarRecDS) (string, []any) {
	return qb.varBuilder.InsertInto(table, rec).Build()
}
