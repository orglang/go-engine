package poolexec

import (
	"database/sql"

	"github.com/huandu/go-sqlbuilder"

	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/implvar"
)

type sqlBuilder struct {
	semBuilder  *sqlbuilder.Struct
	bindBuilder *sqlbuilder.Struct
	execBuilder *sqlbuilder.Struct
	varBuilder  *sqlbuilder.Struct
	snapBuilder *sqlbuilder.Struct
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
	snapBuilder := sqlbuilder.NewStruct(new(execLiabSnapDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{semBuilder, bindBuilder, execBuilder, varBuilder, snapBuilder}
}

func (qb *sqlBuilder) insertRec(rec execRecDS) (string, []any) {
	return qb.execBuilder.InsertInto("pool_execs", rec).Build()
}

func (qb *sqlBuilder) selectSnap(ref implsem.SemRefDS) (string, []any) {
	sems := qb.semBuilder.SelectFrom(implSems)
	svs := qb.varBuilder.SelectFrom(poolStructVars)
	lvs := qb.varBuilder.SelectFrom(poolLinearVars)
	id := sql.Named("id", ref.ImplID)
	return sems.SelectMore(
		sems.BuilderAs(svs.Where(svs.Equal("impl_id", id)), "struct_vars"),
		sems.BuilderAs(lvs.Where(lvs.Equal("impl_id", id)), "linear_vars"),
	).
		Where(sems.Equal("impl_id", id)).
		Build()
}

func (qb *sqlBuilder) selectSnapByQN(implQN string) (string, []any) {
	sb := qb.snapBuilder.SelectFrom(implSems)
	return sb.Join(poolExecs, "exec.impl_id = sem.impl_id").
		Join(poolStructVars, "struct_var.impl_id = sem.impl_id", "struct_var.side = 1", "exec.mode = 1").
		Join(poolLinearVars, "linear_var.impl_id = sem.impl_id", "linear_var.side = 1", "exec.mode = 2").
		Join(implBinds, "bind.impl_id = sem.impl_id").
		Where(sb.Equal("bind.impl_qn", implQN)).
		OrderByDesc("struct_var.impl_rn").OrderByDesc("linear_var.impl_rn").
		Limit(1).
		Build()
}
