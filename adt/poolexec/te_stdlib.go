package poolexec

import (
	"html/template"
	"log/slog"

	"github.com/Masterminds/sprig/v3"

	"orglang/go-engine/lib/te"
)

func newRendererStdlib(l *slog.Logger) (*te.RendererStdlib, error) {
	t := template.New("poolexec").Funcs(sprig.FuncMap())
	return te.NewRendererStdlib(t, l), nil
}
