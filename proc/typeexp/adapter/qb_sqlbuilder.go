package adapter

import (
	"github.com/huandu/go-sqlbuilder"

	proctypeexp "orglang/go-engine/proc/typeexp/core"
)

type sqlBuilder struct {
	stateBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() proctypeexp.QueryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	stateBuilder := sqlbuilder.NewStruct(new(proctypeexp.StateDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{stateBuilder}
}

func (qb *sqlBuilder) InsertRec(rec proctypeexp.StateDS) (string, []any) {
	return qb.stateBuilder.InsertInto(xactExps, rec).
		SQL("ON CONFLICT (exp_vk) DO NOTHING").
		Build()
}

func (qb *sqlBuilder) SelectRecByVK() string {
	top := qb.stateBuilder.SelectFrom(xactExps + " top")
	top.Where(top.Equal("top.exp_vk", top.Var(1)))
	sub := qb.stateBuilder.SelectFrom("exp_tree sup, pool_type_exps sub")
	sub.Where("sub.sup_exp_vk = sup.exp_vk")
	tree := sqlbuilder.PostgreSQL.NewCTEBuilder()
	ub := sqlbuilder.PostgreSQL.NewUnionBuilder()
	tree.WithRecursive(
		sqlbuilder.CTETable("exp_tree").As(ub.UnionAll(top, sub)),
	)
	return tree.Select("*").String()
}
