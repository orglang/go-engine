package adapter

import (
	"io"
	"log/slog"
	"reflect"

	"github.com/labstack/echo/v4"
	sdkprog "github.com/orglang/go-sdk/prog"

	webapi "orglang/go-engine/internal/app/web/api"
	coreprog "orglang/go-engine/internal/prog/core"
)

type controllerEcho struct {
	webapi.ServerInterface
	api coreprog.API
	log *slog.Logger
}

var _ webapi.ServerInterface = (*controllerEcho)(nil)

func newControllerEcho(a coreprog.API, l *slog.Logger) *controllerEcho {
	name := slog.String("name", reflect.TypeFor[controllerEcho]().Name())
	return &controllerEcho{api: a, log: l.With(name)}
}

func (h *controllerEcho) CreateProgram(ctx echo.Context) error {
	text, readErr := io.ReadAll(ctx.Request().Body)
	if readErr != nil {
		h.log.Error("reading failed", slog.Any("err", readErr))
		return readErr
	}
	dto, parseErr := sdkprog.MsgFromText(string(text))
	if parseErr != nil {
		h.log.Error("parsing failed", slog.Any("dto", reflect.TypeFor[sdkprog.Spec]()))
		return parseErr
	}
	spec, convErr := MsgToSpec(dto)
	if convErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return convErr
	}
	if apiErr := h.api.Create(spec); apiErr != nil {
		return apiErr
	}
	return ctx.NoContent(http.StatusCreated)
}
