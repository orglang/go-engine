package semcomm

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
	return qb.semBuilder.InsertInto(commSems, rec).Build()
}

func (qb *sqlBuilder) updateRec(ref SemRefDS) (string, []any) {
	sem := sqlbuilder.PostgreSQL.NewUpdateBuilder()
	sem.Update(commSems)
	sem.Set("comm_rn = comm_rn + 1")
	sem.Where(sem.Equal("comm_id", ref.CommID), sem.Equal("comm_rn", ref.CommRN))
	return sem.Build()
}
