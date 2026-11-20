package adapter

import (
	"github.com/huandu/go-sqlbuilder"

)

type sqlBuilder struct {
	defBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() QueryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	defBuilder := sqlbuilder.NewStruct(new(DefRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{defBuilder}
}

func (qb *sqlBuilder) InsertRec(rec DefRecDS) (string, []any) {
	return qb.defBuilder.InsertInto(typeDefs, rec).Build()
}

func (qb *sqlBuilder) SelectRecByQN() string {
	sb := qb.defBuilder.SelectFrom(typeDefs + "def")
	return sb.Join(descBinds+"bind", "bind.desc_id = def.type_id").
		Where(sb.Equal("bind.desc_qn", sb.Var(1))).
		String()
}
