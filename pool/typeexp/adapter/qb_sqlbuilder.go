package adapter

import (
	"github.com/huandu/go-sqlbuilder"

	pooltypeexp "orglang/go-engine/pool/typeexp/core"
)

type sqlBuilder struct {
	stateBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() pooltypeexp.QueryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	stateBuilder := sqlbuilder.NewStruct(new(pooltypeexp.StateDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{stateBuilder}
}

func (qb *sqlBuilder) InsertRec(rec pooltypeexp.StateDS) (string, []any) {
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
