package termdef

import (
	"github.com/huandu/go-sqlbuilder"
)

type sqlBuilder struct {
	defBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	defBuilder := sqlbuilder.NewStruct(new(defRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{defBuilder}
}

func (qb *sqlBuilder) insertRec(rec defRecDS) (string, []any) {
	return qb.defBuilder.InsertInto(termDefs, rec).Build()
}

func (qb *sqlBuilder) selectRecByQN(qn string) (string, []any) {
	sb := qb.defBuilder.SelectFrom(termDefs + "def")
	return sb.Join(descBinds+"bind", "bind.desc_id = def.term_id").
		Where(sb.Equal("bind.desc_qn", qn)).
		Build()
}
