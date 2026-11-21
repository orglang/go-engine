package adapter

import (
	"github.com/huandu/go-sqlbuilder"

	compvaradapter "orglang/go-engine/internal/adt/compvar/adapter"
	termsemadapter "orglang/go-engine/internal/adt/termsem/adapter"
)

type sqlBuilder struct {
	semBuilder  *sqlbuilder.Struct
	execBuilder *sqlbuilder.Struct
	varBuilder  *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	semBuilder := sqlbuilder.NewStruct(new(termsemadapter.SemRef)).For(sqlbuilder.PostgreSQL)
	execBuilder := sqlbuilder.NewStruct(new(ExecRec)).For(sqlbuilder.PostgreSQL)
	varBuilder := sqlbuilder.NewStruct(new(compvaradapter.VarRec)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{semBuilder, execBuilder, varBuilder}
}

func (qb *sqlBuilder) InsertRec(rec ExecRec) (string, []any) {
	return qb.execBuilder.InsertInto(CompExecs, rec).Build()
}
