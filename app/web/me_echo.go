package web

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"orglang/orglang/avt/msg"

	typedef "orglang/orglang/aat/type/def"
)

// Adapter
type handlerEcho struct {
	api typedef.API
	ssr msg.Renderer
	log *slog.Logger
}

func newHandlerEcho(a typedef.API, r msg.Renderer, l *slog.Logger) *handlerEcho {
	name := slog.String("name", "webHandlerEcho")
	return &handlerEcho{a, r, l.With(name)}
}

func (h *handlerEcho) Home(c echo.Context) error {
	refs, err := h.api.RetreiveRefs()
	if err != nil {
		return err
	}
	html, err := h.ssr.Render("home.html", typedef.MsgFromTypeRefs(refs))
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}
