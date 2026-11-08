package adapter

import (
	"context"
	"errors"
	"log/slog"
	"reflect"

	"github.com/jackc/pgx/v5"

	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/tx"

	"orglang/go-engine/internal/adt/identity"

	commturn "orglang/go-engine/internal/proc/commturn/core"
)

type daoPgx struct {
	log *slog.Logger
}

// for compilation purposes
func newRepo() commturn.Repo {
	return new(daoPgx)
}

func NewDaoPgx(l *slog.Logger) commturn.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{l.With(name)}
}

func (dao *daoPgx) InsertRecs(ctx context.Context, conn tx.Conn, recs ...commturn.TurnRec) error {
	dtos, err := DataFromStepRecs(recs)
	if err != nil {
		dao.log.Error("conversion failed")
		return err
	}
	batch := pgx.Batch{}
	for _, dto := range dtos {
		args := pgx.NamedArgs{
			"kind":    dto.K,
			"comp_id": dto.ExecID,
			"chnl_id": dto.ChnlID,
			"proc_er": dto.ProcER,
		}
		batch.Queue(insertStep, args)
	}
	br := conn.Pgx.SendBatch(ctx, &batch)
	defer func() {
		err = errors.Join(err, br.Close())
	}()
	for _, dto := range dtos {
		_, err = br.Exec()
		if err != nil {
			dao.log.Error("execution failed", slog.String("q", insertStep), slog.Any("dto", dto))
		}
	}
	if err != nil {
		return err
	}
	return nil
}

func (dao *daoPgx) SelectRecs(ctx context.Context, conn tx.Conn, rid identity.ADT) (commturn.TurnRec, error) {
	query := `
		select
			id, kind, pid, vid, spec
		from steps
		WHERE id = $1`
	return dao.execute(ctx, conn, query, rid.String())
}

func (dao *daoPgx) execute(ctx context.Context, conn tx.Conn, query string, arg string) (commturn.TurnRec, error) {
	rows, err := conn.Pgx.Query(ctx, query, arg)
	if err != nil {
		dao.log.Error("query execution failed", slog.String("q", query))
		return nil, err
	}
	defer rows.Close()
	dto, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[commturn.StepRecDS])
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		dao.log.Error("row scanning failed")
		return nil, err
	}
	root, err := dataToStepRec(dto)
	if err != nil {
		dao.log.Error("model conversion failed")
		return nil, err
	}
	dao.log.Log(ctx, lf.LevelTrace, "entity selection succeed", slog.Any("root", root))
	return root, nil
}

const (
	insertStep = `
		insert into pool_comm_turns (
			id, kind, pid, vid, spec
		) values (
			@id, @kind, @pid, @vid, @spec
		)`
)
