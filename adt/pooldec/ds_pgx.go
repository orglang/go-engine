package pooldec

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

func (dao *pgxDAO) InsertRec(source db.Source, rec DecRec) error {
	ds := db.MustConform[db.SourcePgx](source)
	idAttr := slog.Any("id", rec.DescRef.DescID)
	dto, err := DataFromDecRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", idAttr)
		return err
	}
	args := pgx.NamedArgs{
		"desc_id":     dto.DescID,
		"provider_vr": dto.ProviderVR,
		"client_vrs":  dto.ClientVRs,
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
		insert into pool_decs (
			desc_id, provider_vr, client_vrs
		) values (
			@desc_id, @provider_vr, @client_vrs
		)`
)
