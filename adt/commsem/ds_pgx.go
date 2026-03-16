package commsem

import (
	"log/slog"
	"reflect"

	"orglang/go-engine/lib/db"

	"orglang/go-engine/adt/uniqsym"
)

type pgxDAO struct {
	log *slog.Logger
}

func newPgxDAO(l *slog.Logger) *pgxDAO {
	name := slog.String("name", reflect.TypeFor[pgxDAO]().Name())
	return &pgxDAO{l.With(name)}
}

// for compilation purposes
func newRepo() Repo {
	return new(pgxDAO)
}

func (p *pgxDAO) InsertRec(db.Source, SemRec) error {
	panic("unimplemented")
}

func (p *pgxDAO) SelectRefsByQNs(db.Source, []uniqsym.ADT) (map[uniqsym.ADT]SemRef, error) {
	panic("unimplemented")
}

func (p *pgxDAO) TouchRec(db.Source, SemRef) error {
	panic("unimplemented")
}
