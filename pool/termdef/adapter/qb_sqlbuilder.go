package adapter

import (
	"github.com/huandu/go-sqlbuilder"

	termdef "orglang/go-engine/pool/termdef/core"
)

type sqlBuilder struct {
	defBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() termdef.QueryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	defBuilder := sqlbuilder.NewStruct(new(termdef.DefRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{defBuilder}
}

func (qb *sqlBuilder) InsertRec(rec termdef.DefRecDS) (string, []any) {
	return qb.defBuilder.InsertInto(termDefs, rec).Build()
}

func (qb *sqlBuilder) SelectRecByQN(qn string) (string, []any) {
	sb := qb.defBuilder.SelectFrom(termDefs + "def")
	return sb.Join(descBinds+"bind", "bind.desc_id = def.term_id").
		Where(sb.Equal("bind.desc_qn", qn)).
		Build()
}
