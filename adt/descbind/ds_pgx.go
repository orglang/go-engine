package descbind

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

func (dao *pgxDAO) InsertRec(source db.Source, rec BindRec) error {
	ds := db.MustConform[db.SourcePgx](source)
	qnAttr := slog.Any("desc_qn", rec.DescQN)
	dto, err := DataFromRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", qnAttr)
		return err
	}
	args := pgx.NamedArgs{
		"desc_qn": dto.DescQN,
		"desc_id": dto.DescID,
	}
	_, err = ds.Conn.Exec(ds.Ctx, insertRec, args)
	if err != nil {
		dao.log.Error("query execution failed", qnAttr)
		return err
	}
	return nil
}

const (
	insertRec = `
		insert into desc_binds (
			desc_qn, desc_id
		) values (
			@desc_qn, @desc_id
		)`
)
