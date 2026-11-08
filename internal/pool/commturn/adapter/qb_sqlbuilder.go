package adapter

import (
	"github.com/huandu/go-sqlbuilder"

	"orglang/go-engine/adt/termsem"

	commturn "orglang/go-engine/pool/commturn/core"
)

type sqlBuilder struct {
	semBuilder  *sqlbuilder.Struct
	stepBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() commturn.QueryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	semBuilder := sqlbuilder.NewStruct(new(termsem.SemRefDS)).For(sqlbuilder.PostgreSQL)
	stepBuilder := sqlbuilder.NewStruct(new(commturn.TurnRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{semBuilder, stepBuilder}
}

func (qb *sqlBuilder) InsertRec(rec commturn.TurnRecDS) (string, []any) {
	return qb.stepBuilder.InsertInto(commTurns, rec).Build()
}
