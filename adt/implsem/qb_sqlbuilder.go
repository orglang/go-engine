package implsem

import (
	"github.com/huandu/go-sqlbuilder"
)

type sqlBuilder struct {
	semBuilder *sqlbuilder.Struct
}

// for compilation purposes
func newQueryBuikder() queryBuilder {
	return new(sqlBuilder)
}

func newSQLBuilder() *sqlBuilder {
	semBuilder := sqlbuilder.NewStruct(new(semRecDS)).For(sqlbuilder.PostgreSQL)
	return &sqlBuilder{semBuilder}
}

func (qb *sqlBuilder) insertRec(rec semRecDS) (string, []any) {
	return qb.semBuilder.InsertInto(implSems, rec).Build()
}

func (qb *sqlBuilder) updateRec(ref SemRefDS) (string, []any) {
	sem := sqlbuilder.PostgreSQL.NewUpdateBuilder()
	sem.Update(implSems)
	sem.Set("impl_rn = impl_rn + 1")
	sem.Where(sem.Equal("impl_id", ref.ImplID), sem.Equal("impl_rn", ref.ImplRN))
	return sem.Build()
}
