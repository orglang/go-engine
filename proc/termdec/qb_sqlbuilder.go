package termdec

import (
	"github.com/huandu/go-sqlbuilder"
)

type sqlBuilder struct {
	decBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	decBuilder := sqlbuilder.NewStruct(new(decRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{decBuilder}
}

func (qb *sqlBuilder) insertRec(rec decRecDS) (string, []any) {
	return qb.decBuilder.InsertInto(termDecs, rec).Build()
}
