package adapter

import "github.com/huandu/go-sqlbuilder"

type queryBuilder interface {
	InsertRec(DecRecDS) (string, []any)
}

type sqlBuilder struct {
	decBuilder *sqlbuilder.Struct
}

func newQueryBuilder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	decBuilder := sqlbuilder.NewStruct(new(DecRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{decBuilder}
}

func (qb *sqlBuilder) InsertRec(rec DecRecDS) (string, []any) {
	return qb.decBuilder.InsertInto(termDecs, rec).Build()
}
