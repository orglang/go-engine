package poolimpl

import (
	"github.com/huandu/go-sqlbuilder"

	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/implvar"
)

type sqlBuilder struct {
	implRecBuilder *sqlbuilder.Struct
	varRecBuilder  *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuikder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	implJoinBuilder := sqlbuilder.NewStruct(new(implJoinDS)).For(sqlbuilder.PostgreSQL)
	varRecBuilder := sqlbuilder.NewStruct(new(implvar.VarRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{implJoinBuilder, varRecBuilder}
}

func (qb *sqlBuilder) selectRecByRef(ref implsem.SemRefDS) (string, []any) {
	implRec := qb.implRecBuilder.SelectFrom(implSems + "sem")
	structVar := qb.varRecBuilder.SelectFrom(poolStructVars + "var")
	linearVar := qb.varRecBuilder.SelectFrom(poolLinearVars + "var")
	vars := sqlbuilder.PostgreSQL.NewCTEBuilder()
	vars.With(
		sqlbuilder.CTEQuery("struct_vars").As(structVar.Where(structVar.Equal("impl_id", ref.ImplID))),
		sqlbuilder.CTEQuery("linear_vars").As(linearVar.Where(linearVar.Equal("impl_id", ref.ImplID))),
	)
	return implRec.With(vars).
		SelectMore(
			implRec.BuilderAs(sqlbuilder.Buildf(arrayAgg, sqlbuilder.Raw("struct_vars")), "struct_vars"),
			implRec.BuilderAs(sqlbuilder.Buildf(arrayAgg, sqlbuilder.Raw("linear_vars")), "linear_vars"),
		).
		Join(poolExecs+"exec", "exec.impl_id = sem.impl_id").
		Where(implRec.Equal("sem.impl_id", ref.ImplID)).
		Build()
}

const (
	arrayAgg = "SELECT array_agg(row(r.*)) FROM %s r"
)
