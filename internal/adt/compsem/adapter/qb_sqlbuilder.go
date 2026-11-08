package adapter

import (
	"github.com/huandu/go-sqlbuilder"

	"orglang/go-engine/internal/adt/compsem"
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
	semBuilder := sqlbuilder.NewStruct(new(compsem.SemRefDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{tableName, semBuilder}
}

func (qb *sqlBuilder) updateRef(ref compsem.SemRefDS) (string, []any) {
	sem := sqlbuilder.PostgreSQL.NewUpdateBuilder()
	sem.Update(qb.tableName)
	sem.Set("comp_rn = comp_rn + 1")
	sem.Where(sem.Equal("comp_id", ref.CompID), sem.Equal("comp_rn", ref.CompRN))
	return sem.Build()
}
