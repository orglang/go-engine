package typesem

import (
	"github.com/huandu/go-sqlbuilder"
)

type sqlBuilder struct {
	typeTable  string
	descTable  string
	refBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuilder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder(typeTable, descTable string) *sqlBuilder {
	refBuilder := sqlbuilder.NewStruct(new(SemRefDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{typeTable, descTable, refBuilder}
}

func (qb *sqlBuilder) updateRef(ref SemRefDS) (string, []any) {
	sem := sqlbuilder.PostgreSQL.NewUpdateBuilder()
	sem.Update(qb.typeTable)
	sem.Set("comp_rn = comp_rn + 1")
	sem.Where(sem.Equal("comp_id", ref.TypeID), sem.Equal("comp_rn", ref.TypeRN))
	return sem.Build()
}

func (qb *sqlBuilder) selectRefByQN() string {
	sb := qb.refBuilder.SelectFrom(qb.typeTable + "t")
	return sb.Join(qb.descTable+"d", "d.desc_id = t.type_id").
		Where(sb.Equal("d.desc_qn", sb.Var(1))).
		String()
}
