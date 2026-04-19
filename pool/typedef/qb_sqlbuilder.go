package typedef

import (
	"github.com/huandu/go-sqlbuilder"
)

type sqlBuilder struct {
	defBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuikder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	defBuilder := sqlbuilder.NewStruct(new(defRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{defBuilder}
}

func (qb *sqlBuilder) insertRec(rec defRecDS) (string, []any) {
	return qb.defBuilder.InsertInto(typeDefs, rec).Build()
}

func (qb *sqlBuilder) selectRecByQN() string {
	sb := qb.defBuilder.SelectFrom(typeDefs + "def")
	return sb.Join(descBinds+"bind", "bind.desc_id = def.type_id").
		Where(sb.Equal("bind.desc_qn", sb.Var(1))).
		String()
}
