package adapter

import (
	"github.com/huandu/go-sqlbuilder"

	"orglang/go-engine/internal/adt/descsem"
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
	semBuilder := sqlbuilder.NewStruct(new(descsem.SemRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{tableName, semBuilder}
}

func (qb *sqlBuilder) insertRec(rec descsem.SemRecDS) (string, []any) {
	return qb.semBuilder.InsertInto(qb.tableName, rec).Build()
}
