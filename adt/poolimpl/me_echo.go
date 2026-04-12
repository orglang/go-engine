package poolimpl

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	sdk "github.com/orglang/go-sdk/adt/poolstep"

	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/poolstep"
)

// Server-side primary adapter
type echoController struct {
	api API
	log *slog.Logger
}

func newEchoController(api API, log *slog.Logger) *echoController {
	name := slog.String("name", reflect.TypeFor[echoController]().Name())
	return &echoController{api, log.With(name)}
}

func cfgEchoController(server *echo.Echo, controller *echoController) error {
	server.POST("/api/v1/pools/execs/steps", controller.PostSpec2)
	server.POST("/api/v1/pools/execs/spawns", controller.PostSpec3)
	return nil
}

func (c *echoController) PostSpec2(ctx echo.Context) error {
	var dto sdk.StepSpec
	bindErr := ctx.Bind(&dto)
	if bindErr != nil {
		c.log.Error("binding failed", slog.Any("dto", reflect.TypeFor[sdk.StepSpec]()))
		return bindErr
	}
	validErr := dto.Validate()
	if validErr != nil {
		c.log.Error("validation failed", slog.Any("dto", dto))
		return validErr
	}
	spec, convErr := poolstep.MsgToStepSpec(dto)
	if convErr != nil {
		c.log.Error("conversion failed", slog.Any("dto", dto))
		return convErr
	}
	apiErr := c.api.Take(spec)
	if apiErr != nil {
		return apiErr
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (c *echoController) PostSpec3(ctx echo.Context) error {
	var dto sdk.StepSpec
	bindErr := ctx.Bind(&dto)
	if bindErr != nil {
		c.log.Error("binding failed", slog.Any("dto", reflect.TypeFor[sdk.StepSpec]()))
		return bindErr
	}
	validateErr := dto.Validate()
	if validateErr != nil {
		c.log.Error("validation failed", slog.Any("dto", dto))
		return validateErr
	}
	spec, convErr := poolstep.MsgToStepSpec(dto)
	if convErr != nil {
		c.log.Error("conversion failed", slog.Any("dto", dto))
		return convErr
	}
	ref, apiErr := c.api.Spawn(spec)
	if apiErr != nil {
		return apiErr
	}
	return ctx.JSON(http.StatusCreated, implsem.MsgFromRef(ref))
}
