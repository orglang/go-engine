package procexec

import (
	"github.com/huandu/go-sqlbuilder"

	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/implvar"
)

type sqlBuilder struct {
	semBuilder  *sqlbuilder.Struct
	bindBuilder *sqlbuilder.Struct
	execBuilder *sqlbuilder.Struct
	varBuilder  *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuikder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	semBuilder := sqlbuilder.NewStruct(new(implsem.SemRefDS)).For(sqlbuilder.PostgreSQL)
	bindBuilder := sqlbuilder.NewStruct(new(implsem.SemBindDS)).For(sqlbuilder.PostgreSQL)
	execBuilder := sqlbuilder.NewStruct(new(execRecDS)).For(sqlbuilder.PostgreSQL)
	varBuilder := sqlbuilder.NewStruct(new(implvar.VarRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{semBuilder, bindBuilder, execBuilder, varBuilder}
}

func (qb *sqlBuilder) insertRec(rec execRecDS) (string, []any) {
	return qb.execBuilder.InsertInto(procExecs, rec).Build()
}
