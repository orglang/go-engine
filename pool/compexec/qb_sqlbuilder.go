package compexec

import (
	"github.com/huandu/go-sqlbuilder"

	"orglang/go-engine/adt/compvar"
	"orglang/go-engine/adt/semcomp"
)

type sqlBuilder struct {
	implJoinBuilder *sqlbuilder.Struct
	varRecBuilder   *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuikder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	implJoinBuilder := sqlbuilder.NewStruct(new(implJoinDS)).For(sqlbuilder.PostgreSQL)
	varRecBuilder := sqlbuilder.NewStruct(new(compvar.VarRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{implJoinBuilder, varRecBuilder}
}

func (qb *sqlBuilder) selectRecByRef(ref semcomp.CompRefDS) (string, []any) {
	poolImpl := qb.implJoinBuilder.SelectFrom(implSems + "sem")
	structVar := qb.varRecBuilder.SelectFrom(poolStructVars + "var")
	linearVar := qb.varRecBuilder.SelectFrom(poolLinearVars + "var")
	vars := sqlbuilder.PostgreSQL.NewCTEBuilder()
	vars.With(
		sqlbuilder.CTEQuery("struct_vars").As(structVar.Where(structVar.Equal("impl_id", ref.CompID))),
		sqlbuilder.CTEQuery("linear_vars").As(linearVar.Where(linearVar.Equal("impl_id", ref.CompID))),
	)
	return poolImpl.With(vars).
		SelectMore(
			poolImpl.BuilderAs(sqlbuilder.Buildf(arrayAgg, sqlbuilder.Raw("struct_vars")), "struct_vars"),
			poolImpl.BuilderAs(sqlbuilder.Buildf(arrayAgg, sqlbuilder.Raw("linear_vars")), "linear_vars"),
		).
		Join(poolExecs+"exec", "exec.impl_id = sem.impl_id").
		Where(poolImpl.Equal("sem.impl_id", ref.CompID)).
		Build()
}

const (
	arrayAgg = "SELECT array_agg(row(r.*)) FROM %s r"
)
