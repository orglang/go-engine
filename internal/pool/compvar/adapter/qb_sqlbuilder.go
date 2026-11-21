package adapter

import (
	"github.com/huandu/go-sqlbuilder"

	compvaradapter "orglang/go-engine/internal/adt/compvar/adapter"
)

type sqlBuilder struct {
	varBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	varBuilder := sqlbuilder.NewStruct(new(compvaradapter.VarRec)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{varBuilder}
}

func (qb *sqlBuilder) InsertRec(table string, rec compvaradapter.VarRec) (string, []any) {
	return qb.varBuilder.InsertInto(table, rec).Build()
}
