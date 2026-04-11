package poolexec

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	"github.com/orglang/go-sdk/adt/poolexec"

	"orglang/go-engine/lib/te"

	"orglang/go-engine/adt/implsem"
)

// Server-side primary adapter
type echoController struct {
	api API
	ssr te.Renderer
	log *slog.Logger
}

func newEchoController(api API, ssr te.Renderer, log *slog.Logger) *echoController {
	name := slog.String("name", reflect.TypeFor[echoController]().Name())
	return &echoController{api, ssr, log.With(name)}
}

func cfgEchoController(server *echo.Echo, controller *echoController) error {
	server.POST("/api/v1/pools/execs", controller.PostSpec)
	return nil
}

func (c *echoController) PostSpec(ctx echo.Context) error {
	var dto poolexec.ExecSpec
	bindErr := ctx.Bind(&dto)
	if bindErr != nil {
		c.log.Error("binding failed", slog.Any("dto", reflect.TypeOf(dto)))
		return bindErr
	}
	validErr := dto.Validate()
	if validErr != nil {
		c.log.Error("validation failed", slog.Any("dto", dto))
		return validErr
	}
	spec, convErr := MsgToExecSpec(dto)
	if convErr != nil {
		c.log.Error("conversion failed", slog.Any("dto", dto))
		return convErr
	}
	ref, apiErr := c.api.Run(spec)
	if apiErr != nil {
		return apiErr
	}
	return ctx.JSON(http.StatusCreated, implsem.MsgFromRef(ref))
}
