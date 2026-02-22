package poolvar

import (
	"errors"
	"log/slog"
	"orglang/go-engine/adt/implvar"
	"orglang/go-engine/lib/db"
	"orglang/go-engine/lib/lf"
	"reflect"

	"github.com/jackc/pgx/v5"
)

type pgxDAO struct {
	log *slog.Logger
}

func newPgxDAO(log *slog.Logger) *pgxDAO {
	name := slog.String("name", reflect.TypeFor[pgxDAO]().Name())
	return &pgxDAO{log.With(name)}
}

// for compilation purposes
func newRepo() Repo {
	return new(pgxDAO)
}

func (dao *pgxDAO) InsertRecs(source db.Source, recs []implvar.VarRec) (err error) {
	ds := db.MustConform[db.SourcePgx](source)
	dtos := implvar.DataFromVarRecs(recs)
	batch := pgx.Batch{}
	for _, dto := range dtos {
		args := pgx.NamedArgs{
			"chnl_id": dto.ChnlID,
			"chnl_ph": dto.ChnlPH,
			"exp_vk":  dto.ExpVK,
			"impl_id": dto.ImplID,
			"impl_rn": dto.ImplRN,
		}
		batch.Queue(insertRec, args)
	}
	br := ds.Conn.SendBatch(ds.Ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	for _, dto := range dtos {
		_, readErr := br.Exec()
		if readErr != nil {
			dao.log.Error("query execution failed", slog.Any("dto", dto))
			return readErr
		}
	}
	dao.log.Log(ds.Ctx, lf.LevelTrace, "insertion succeed", slog.Any("dtos", dtos))
	return nil
}

const (
	insertRec = `
		insert into pool_vars (
			chnl_id, chnl_ph, exp_vk, impl_id, impl_rn
		) values (
			@chnl_id, @chnl_ph, @exp_vk, @impl_id, @impl_rn
		)`
)
