package synonym

import (
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/lib/db"
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

func (dao *pgxDAO) InsertRec(source db.Source, rec Rec) error {
	ds := db.MustConform[db.SourcePgx](source)
	idAttr := slog.Any("decID", rec.SynVK)
	dto, err := DataFromRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", idAttr)
		return err
	}
	args := pgx.NamedArgs{
		"syn_vk": dto.SynVK,
		"syn_qn": dto.SynQN,
	}
	_, err = ds.Conn.Exec(ds.Ctx, insertRec, args)
	if err != nil {
		dao.log.Error("query execution failed", idAttr)
		return err
	}
	return nil
}

const (
	insertRec = `
		insert into synonyms (
			syn_vk, syn_qn
		) values (
			@syn_vk, @syn_qn
		)`
)
