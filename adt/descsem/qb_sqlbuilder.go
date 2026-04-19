package descsem

import (
	"github.com/huandu/go-sqlbuilder"
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
	semBuilder := sqlbuilder.NewStruct(new(SemRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{tableName, semBuilder}
}

func (qb *sqlBuilder) insertRec(rec SemRecDS) (string, []any) {
	return qb.semBuilder.InsertInto(qb.tableName, rec).Build()
}
