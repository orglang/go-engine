package adapter

import (
	"github.com/huandu/go-sqlbuilder"

	implsem "orglang/go-engine/internal/adt/implsem"
)

type sqlBuilder struct {
	tableName  string
	semBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder(tableName string) *sqlBuilder {
	semBuilder := sqlbuilder.NewStruct(new(implsem.SemRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{tableName, semBuilder}
}

func (qb *sqlBuilder) insertRec(rec implsem.SemRecDS) (string, []any) {
	return qb.semBuilder.InsertInto(qb.tableName, rec).Build()
}
