package adapter

import (
	"log/slog"
	"reflect"

	pooltermexp "orglang/go-engine/pool/termexp/core"
)

// Adapter
type daoPgx struct {
	log *slog.Logger
}

func NewDaoPgx(l *slog.Logger) pooltermexp.Repo {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{l.With(name)}
}
