package compexec

import (
	"github.com/huandu/go-sqlbuilder"

	"orglang/go-engine/adt/compvar"
	"orglang/go-engine/adt/termsem"
)

type sqlBuilder struct {
	semBuilder  *sqlbuilder.Struct
	execBuilder *sqlbuilder.Struct
	varBuilder  *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuikder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	semBuilder := sqlbuilder.NewStruct(new(termsem.SemRefDS)).For(sqlbuilder.PostgreSQL)
	execBuilder := sqlbuilder.NewStruct(new(execRecDS)).For(sqlbuilder.PostgreSQL)
	varBuilder := sqlbuilder.NewStruct(new(compvar.VarRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{semBuilder, execBuilder, varBuilder}
}

func (qb *sqlBuilder) insertRec(rec execRecDS) (string, []any) {
	return qb.execBuilder.InsertInto(compExecs, rec).Build()
}
