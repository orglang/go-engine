package adapter

import (
	"github.com/huandu/go-sqlbuilder"

	typedef "orglang/go-engine/proc/typedef/core"
)

type sqlBuilder struct {
	defBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() typedef.QueryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	defBuilder := sqlbuilder.NewStruct(new(typedef.DefRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{defBuilder}
}

func (qb *sqlBuilder) InsertRec(rec typedef.DefRecDS) (string, []any) {
	return qb.defBuilder.InsertInto(typeDefs, rec).Build()
}

func (qb *sqlBuilder) SelectRecByQN() string {
	sb := qb.defBuilder.SelectFrom(typeDefs + "def")
	return sb.Join(descBinds+"bind", "bind.desc_id = def.type_id").
		Where(sb.Equal("bind.desc_qn", sb.Var(1))).
		String()
}
