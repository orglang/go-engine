package adapter

import "github.com/huandu/go-sqlbuilder"

type queryBuilder interface {
	InsertRec(DefRecDS) (string, []any)
	SelectRecByQN(string) (string, []any)
}

type sqlBuilder struct {
	defBuilder *sqlbuilder.Struct
}

func newQueryBuilder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	defBuilder := sqlbuilder.NewStruct(new(DefRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{defBuilder}
}

func (qb *sqlBuilder) InsertRec(rec DefRecDS) (string, []any) {
	return qb.defBuilder.InsertInto(termDefs, rec).Build()
}

func (qb *sqlBuilder) SelectRecByQN(qn string) (string, []any) {
	sb := qb.defBuilder.SelectFrom(termDefs + "def")
	return sb.Join(descBinds+"bind", "bind.desc_id = def.term_id").
		Where(sb.Equal("bind.desc_qn", qn)).
		Build()
}
