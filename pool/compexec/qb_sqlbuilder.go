package compexec

import (
	"github.com/huandu/go-sqlbuilder"

	"orglang/go-engine/adt/compsem"
	"orglang/go-engine/adt/compvar"
)

type sqlBuilder struct {
	compExecBuilder *sqlbuilder.Struct
	snapBuilder     *sqlbuilder.Struct
	execJoinBuilder *sqlbuilder.Struct
	compVarBuilder  *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuikder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	compExecBuilder := sqlbuilder.NewStruct(new(execRecDS)).For(sqlbuilder.PostgreSQL)
	snapBuilder := sqlbuilder.NewStruct(new(execSnapDS)).For(sqlbuilder.PostgreSQL)
	execJoinBuilder := sqlbuilder.NewStruct(new(execJoinDS)).For(sqlbuilder.PostgreSQL)
	compVarBuilder := sqlbuilder.NewStruct(new(compvar.VarRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{compExecBuilder, snapBuilder, execJoinBuilder, compVarBuilder}
}

func (qb *sqlBuilder) insertRec(rec execRecDS) (string, []any) {
	return qb.compExecBuilder.InsertInto("pool_comp_execs", rec).Build()
}

func (qb *sqlBuilder) selectRecByRef(ref compsem.SemRefDS) (string, []any) {
	compExec := qb.execJoinBuilder.SelectFrom(compExecs + "sem")
	structVar := qb.compVarBuilder.SelectFrom(poolStructVars + "var")
	linearVar := qb.compVarBuilder.SelectFrom(poolLinearVars + "var")
	vars := sqlbuilder.PostgreSQL.NewCTEBuilder()
	vars.With(
		sqlbuilder.CTEQuery("struct_vars").As(structVar.Where(structVar.Equal("comp_id", ref.CompID))),
		sqlbuilder.CTEQuery("linear_vars").As(linearVar.Where(linearVar.Equal("comp_id", ref.CompID))),
	)
	return compExec.With(vars).
		SelectMore(
			compExec.BuilderAs(sqlbuilder.Buildf(arrayAgg, sqlbuilder.Raw("struct_vars")), "struct_vars"),
			compExec.BuilderAs(sqlbuilder.Buildf(arrayAgg, sqlbuilder.Raw("linear_vars")), "linear_vars"),
		).
		Join(compExecs+"exec", "exec.comp_id = sem.comp_id").
		Where(compExec.Equal("sem.comp_id", ref.CompID)).
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
