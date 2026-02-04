package web

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	"orglang/go-engine/adt/typedef"
	"orglang/go-engine/lib/te"
)

// Adapter
type echoController struct {
	api typedef.API
	ssr te.Renderer
	log *slog.Logger
}

func newEchoController(a typedef.API, r te.Renderer, l *slog.Logger) *echoController {
	name := slog.String("name", reflect.TypeFor[echoController]().Name())
	return &echoController{a, r, l.With(name)}
}

func cfgEchoController(e *echo.Echo, h *echoController) {
	e.GET("/", h.Home)
}

func (h *echoController) Home(c echo.Context) error {
	refs, err := h.api.RetreiveRefs()
	if err != nil {
		return err
	}
	html, err := h.ssr.Render("home.html", typedef.MsgFromDefRefs(refs))
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}
