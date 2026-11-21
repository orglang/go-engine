package adapter

import (
	"log/slog"
	"reflect"

	)

// Adapter
type daoPgx struct {
	log *slog.Logger
}

func NewDaoPgx(l *slog.Logger) *daoPgx {
	name := slog.String("name", reflect.TypeFor[daoPgx]().Name())
	return &daoPgx{l.With(name)}
}
