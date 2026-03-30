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

func (dao *pgxDAO) AddRec(source db.Source, rec DecRec) error {
	ds := db.MustConform[db.SourcePgx](source)
	refAttr := slog.Any("ref", rec.DescRef)
	dto, convErr := DataFromDecRec(rec)
	if convErr != nil {
		dao.log.Error("model conversion failed", refAttr)
		return convErr
	}
	args := pgx.NamedArgs{
		"desc_id":    dto.DescID,
		"liab_var":   dto.LiabVar,
		"asset_vars": dto.AssetVars,
	}
	_, execErr := ds.Conn.Exec(ds.Ctx, insertRec, args)
	if execErr != nil {
		dao.log.Error("query execution failed", refAttr)
		return execErr
	}
	return nil
}

const (
	insertRec = `
		insert into pool_decs (
			desc_id, liab_var, asset_vars
		) values (
			@desc_id, @liab_var, @asset_vars
		)`
)
