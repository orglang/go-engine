package procdec

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	"orglang/orglang/lib/te"

	"orglang/orglang/adt/identity"
)

// Server-side primary adapter
type handlerEcho struct {
	api API
	ssr te.Renderer
	log *slog.Logger
}

func newHandlerEcho(a API, r te.Renderer, l *slog.Logger) *handlerEcho {
	name := slog.String("name", reflect.TypeFor[handlerEcho]().Name())
	return &handlerEcho{a, r, l.With(name)}
}

func cfgHandlerEcho(e *echo.Echo, h *handlerEcho) error {
	e.POST("/api/v1/declarations", h.PostOne)
	e.GET("/api/v1/declarations/:id", h.GetOne)
	return nil
}

func (h *handlerEcho) PostOne(c echo.Context) error {
	var dto DecSpecME
	err := c.Bind(&dto)
	if err != nil {
		h.log.Error("dto binding failed", slog.Any("reason", err))
		return err
	}
	err = dto.Validate()
	if err != nil {
		h.log.Error("dto validation failed", slog.Any("reason", err), slog.Any("dto", dto))
		return err
	}
	spec, err := MsgToDecSpec(dto)
	if err != nil {
		h.log.Error("dto conversion failed", slog.Any("reason", err), slog.Any("dto", dto))
		return err
	}
	snap, err := h.api.Create(spec)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, MsgFromDecSnap(snap))
}

func (h *handlerEcho) GetOne(c echo.Context) error {
	var dto IdentME
	err := c.Bind(&dto)
	if err != nil {
		return err
	}
	id, err := identity.ConvertFromString(dto.DecID)
	if err != nil {
		return err
	}
	snap, err := h.api.Retrieve(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, MsgFromDecSnap(snap))
}
