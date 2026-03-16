package poolexec

import (
	"database/sql"

	"github.com/huandu/go-sqlbuilder"

	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/implvar"
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
	semBuilder := sqlbuilder.NewStruct(new(implsem.SemRefDS)).For(sqlbuilder.PostgreSQL)
	execBuilder := sqlbuilder.NewStruct(new(execRecDS)).For(sqlbuilder.PostgreSQL)
	varBuilder := sqlbuilder.NewStruct(new(implvar.VarRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{semBuilder, execBuilder, varBuilder}
}

func (qb *sqlBuilder) insertRec(rec execRecDS) (string, []any) {
	return qb.execBuilder.InsertInto(execs, rec).Build()
}

func (qb *sqlBuilder) selectSnap(ref implsem.SemRefDS) (string, []any) {
	sems := qb.semBuilder.SelectFrom(sems)
	svs := qb.varBuilder.SelectFrom(structVars)
	lvs := qb.varBuilder.SelectFrom(linearVars)
	implID := sql.Named("id", ref.ImplID)
	return sems.SelectMore(
		sems.BuilderAs(svs.Where(svs.Equal("impl_id", implID)), "struct_vars"),
		sems.BuilderAs(lvs.Where(lvs.Equal("impl_id", implID)), "linear_vars"),
	).
		Where(sems.Equal("impl_id", implID)).
		Build()
}
