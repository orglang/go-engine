package xactdef

import (
	"github.com/huandu/go-sqlbuilder"

	"orglang/go-engine/adt/descsem"
)

type sqlBuilder struct {
	semBuilder  *sqlbuilder.Struct
	bindBuilder *sqlbuilder.Struct
	defBuilder  *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuikder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	semBuilder := sqlbuilder.NewStruct(new(descsem.SemRefDS)).For(sqlbuilder.PostgreSQL)
	bindBuilder := sqlbuilder.NewStruct(new(descsem.SemBindDS)).For(sqlbuilder.PostgreSQL)
	defBuilder := sqlbuilder.NewStruct(new(defRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{semBuilder, bindBuilder, defBuilder}
}

func (qb *sqlBuilder) selectRecByQN() string {
	sb := qb.defBuilder.SelectFrom(xactDefs)
	return sb.Join(descBinds, "bind.desc_id = def.desc_id").
		Where(sb.Equal("bind.desc_qn", sb.Var(1))).
		String()
}
