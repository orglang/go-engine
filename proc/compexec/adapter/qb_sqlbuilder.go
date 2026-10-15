package adapter

import (
	"github.com/huandu/go-sqlbuilder"

	"orglang/go-engine/adt/compvar"
	"orglang/go-engine/adt/termsem"

	compexec "orglang/go-engine/proc/compexec/core"
)

type sqlBuilder struct {
	semBuilder  *sqlbuilder.Struct
	execBuilder *sqlbuilder.Struct
	varBuilder  *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() compexec.QueryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	semBuilder := sqlbuilder.NewStruct(new(termsem.SemRefDS)).For(sqlbuilder.PostgreSQL)
	execBuilder := sqlbuilder.NewStruct(new(compexec.ExecRecDS)).For(sqlbuilder.PostgreSQL)
	varBuilder := sqlbuilder.NewStruct(new(compvar.VarRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{semBuilder, execBuilder, varBuilder}
}

func (qb *sqlBuilder) InsertRec(rec compexec.ExecRecDS) (string, []any) {
	return qb.execBuilder.InsertInto(compexec.CompExecs, rec).Build()
}
