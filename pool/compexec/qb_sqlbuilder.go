package compexec

import (
	"github.com/huandu/go-sqlbuilder"

	"orglang/go-engine/adt/compsem"
	"orglang/go-engine/adt/compvar"
)

type sqlBuilder struct {
	recBuilder     *sqlbuilder.Struct
	snapBuilder    *sqlbuilder.Struct
	compVarBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	recBuilder := sqlbuilder.NewStruct(new(execRec)).For(sqlbuilder.PostgreSQL)
	snapBuilder := sqlbuilder.NewStruct(new(execSnap1)).For(sqlbuilder.PostgreSQL)
	compVarBuilder := sqlbuilder.NewStruct(new(compvar.VarRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{recBuilder, snapBuilder, compVarBuilder}
}

func (qb *sqlBuilder) insertRec(rec execRec) (string, []any) {
	return qb.recBuilder.InsertInto(compExecs, rec).Build()
}

func (qb *sqlBuilder) selectRecByRef(ref compsem.SemRefDS) (string, []any) {
	exec := qb.recBuilder.SelectFrom(compExecs + "exec")
	structVar := qb.compVarBuilder.SelectFrom(poolStructVars + "var")
	linearVar := qb.compVarBuilder.SelectFrom(poolLinearVars + "var")
	vars := sqlbuilder.PostgreSQL.NewCTEBuilder()
	vars.With(
		sqlbuilder.CTEQuery("struct_vars").As(structVar.Where(structVar.Equal("comp_id", ref.CompID))),
		sqlbuilder.CTEQuery("linear_vars").As(linearVar.Where(linearVar.Equal("comp_id", ref.CompID))),
	)
	return exec.With(vars).
		SelectMore(
			exec.BuilderAs(sqlbuilder.Buildf(arrayAgg, sqlbuilder.Raw("struct_vars")), "struct_vars"),
			exec.BuilderAs(sqlbuilder.Buildf(arrayAgg, sqlbuilder.Raw("linear_vars")), "linear_vars"),
		).
		Where(exec.Equal("exec.comp_id", ref.CompID)).
		Build()
}

func (qb *sqlBuilder) selectSnapByQN(qn string) (string, []any) {
	sb := qb.snapBuilder.SelectFrom(compExecs)
	return sb.Join(compExecs, "exec.comp_id = sem.comp_id").
		Join(poolStructVars, "struct_var.comp_id = sem.comp_id", "struct_var.side = 1", "exec.mode = 1").
		Join(poolLinearVars, "linear_var.comp_id = sem.comp_id", "linear_var.side = 1", "exec.mode = 2").
		Join(implBinds, "bind.impl_id = sem.comp_id").
		Where(sb.Equal("bind.impl_qn", qn)).
		OrderByDesc("struct_var.comp_rn").OrderByDesc("linear_var.comp_rn").
		Limit(1).
		Build()
}

const (
	arrayAgg = "SELECT array_agg(row(r.*)) FROM %s r"
)
