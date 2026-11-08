package adapter

import (
	"github.com/huandu/go-sqlbuilder"

	termdec "orglang/go-engine/proc/termdec/core"
)

type sqlBuilder struct {
	decBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() termdec.QueryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	decBuilder := sqlbuilder.NewStruct(new(termdec.DecRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{decBuilder}
}

func (qb *sqlBuilder) InsertRec(rec termdec.DecRecDS) (string, []any) {
	return qb.decBuilder.InsertInto(termDecs, rec).Build()
}
