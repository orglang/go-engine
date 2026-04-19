package termdef

import (
	"github.com/huandu/go-sqlbuilder"
)

type sqlBuilder struct {
	recBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuikder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	recBuilder := sqlbuilder.NewStruct(new(decRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{recBuilder}
}

func (qb *sqlBuilder) insertRec(rec decRecDS) (string, []any) {
	return qb.recBuilder.InsertInto(poolDecs, rec).Build()
}

func (qb *sqlBuilder) selectRecByQN(qn string) (string, []any) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	return sb.Select("dec.desc_id", "dec.liab_var", "dec.asset_vars").
		From(poolDecs+" dec").
		Join(sigBinds+" bind", "bind.desc_id = dec.desc_id").
		Where(sb.Equal("bind.desc_qn", qn)).
		Build()
}
