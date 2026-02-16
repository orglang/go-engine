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
func newRepo() repo {
	return new(pgxDAO)
}

func (dao *pgxDAO) InsertRec(source db.Source, rec DecRec) error {
	ds := db.MustConform[db.SourcePgx](source)
	refAttr := slog.Any("decRef", rec.DecRef)
	dto, err := DataFromDecRec(rec)
	if err != nil {
		dao.log.Error("model conversion failed", refAttr)
		return err
	}
	args := pgx.NamedArgs{
		"dec_id":      dto.ID,
		"dec_rn":      dto.RN,
		"syn_vk":      dto.SynVK,
		"provider_br": dto.ProviderBR,
		"client_brs":  dto.ClientBRs,
	}
	_, err = ds.Conn.Exec(ds.Ctx, insertRec, args)
	if err != nil {
		dao.log.Error("query execution failed", refAttr)
		return err
	}
	return nil
}

const (
	insertRec = `
		insert into pool_decs (
			dec_id, dec_rn, syn_vk, provider_br, client_brs
		) values (
			@dec_id, @dec_rn, @syn_vk, @provider_br, @client_brs
		)`
)
