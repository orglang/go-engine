package adapter

import (
	"log/slog"
	"reflect"

	termexp "orglang/go-engine/pool/termexp/core"
)

// Adapter
type daoPgx struct {
	log *slog.Logger
}

func NewDaoPgx(l *slog.Logger) termexp.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{l.With(name)}
}
