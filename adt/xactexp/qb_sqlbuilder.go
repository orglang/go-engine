package xactexp

import (
	"github.com/huandu/go-sqlbuilder"
)

type sqlBuilder struct {
	stateBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuikder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	stateBuilder := sqlbuilder.NewStruct(new(stateDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{stateBuilder}
}

func (qb *sqlBuilder) insertRec(rec stateDS) (string, []any) {
	return qb.stateBuilder.InsertInto(xactExps, rec).
		SQL("ON CONFLICT (exp_vk) DO NOTHING").
		Build()
}

func (qb *sqlBuilder) selectRecByVK() string {
	ub := sqlbuilder.PostgreSQL.NewUnionBuilder()
	// top := qb.stateBuilder.SelectFrom("xact_exps top")
	top := sqlbuilder.PostgreSQL.NewSelectBuilder()
	top.Select("top.exp_vk", "top.sup_exp_vk", "top.kind", "top.spec")
	top.From(xactExps + " top")
	top.Where(top.Equal("top.exp_vk", top.Var(1)))
	// sub := qb.stateBuilder.SelectFrom("xact_exps sub, exp_tree sup")
	sub := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sub.Select("sub.exp_vk", "sub.sup_exp_vk", "sub.kind", "sub.spec")
	sub.From(xactExps+" sub", "exp_tree sup")
	sub.Where("sub.sup_exp_vk = sup.exp_vk")
	tree := sqlbuilder.PostgreSQL.NewCTEBuilder()
	tree.WithRecursive(
		sqlbuilder.CTETable("exp_tree").As(ub.UnionAll(top, sub)),
	)
	return tree.Select("*").String()
	// selectByVK := `
	// with recursive exp_tree AS (
	// 	select top.*
	// 	from xact_exps top
	// 	where exp_vk = $1
	// 	union all
	// 	select sub.*
	// 	from xact_exps sub, exp_tree sup
	// 	where sub.sup_exp_vk = sup.exp_vk
	// )
	// select * from exp_tree`
	// return selectByVK, []any{vk}
}
