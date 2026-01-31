package procdef

import (
	"log/slog"
	"orglang/go-runtime/lib/db"
	"reflect"
)

// Adapter
type pgxDAO struct {
	log *slog.Logger
}

// for compilation purposes
func newRepo() Repo {
	return &pgxDAO{}
}

func newPgxDAO(l *slog.Logger) *pgxDAO {
	name := slog.String("name", reflect.TypeFor[pgxDAO]().Name())
	return &pgxDAO{l.With(name)}
}

func (dao *pgxDAO) InsertProc(db.Source, DefRec) error {
	return nil
}
