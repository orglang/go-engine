package poolexec

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	sdk "github.com/orglang/go-sdk/adt/implsem"
	"github.com/orglang/go-sdk/adt/poolexec"
	poolstep1 "github.com/orglang/go-sdk/adt/poolstep"

	"orglang/go-engine/lib/te"

	"orglang/go-engine/adt/implsem"
	"orglang/go-engine/adt/poolstep"
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

func cfgEchoController(e *echo.Echo, h *echoController) error {
	e.POST("/api/v1/pools/execs", h.PostSpec)
	e.GET("/api/v1/pools/execs/:id", h.GetSnap)
	e.POST("/api/v1/pools/execs/steps", h.PostSpec2)
	e.POST("/api/v1/pools/execs/spawns", h.PostSpec3)
	return nil
}

func (c *echoController) PostSpec(ctx echo.Context) error {
	var dto poolexec.ExecSpec
	bindErr := ctx.Bind(&dto)
	if bindErr != nil {
		c.log.Error("binding failed", slog.Any("dto", reflect.TypeOf(dto)))
		return bindErr
	}
	validateErr := dto.Validate()
	if validateErr != nil {
		c.log.Error("validation failed", slog.Any("dto", dto))
		return validateErr
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

func (c *echoController) GetSnap(ctx echo.Context) error {
	var dto sdk.SemRef
	bindErr := ctx.Bind(&dto)
	if bindErr != nil {
		return bindErr
	}
	ref, convErr := implsem.MsgToRef(dto)
	if convErr != nil {
		return convErr
	}
	snap, apiErr := c.api.RetrieveSnap(ref)
	if apiErr != nil {
		return apiErr
	}
	return ctx.JSON(http.StatusOK, MsgFromExecSnap(snap))
}

func (c *echoController) PostSpec2(ctx echo.Context) error {
	var dto poolstep1.StepSpec
	bindErr := ctx.Bind(&dto)
	if bindErr != nil {
		c.log.Error("binding failed", slog.Any("dto", reflect.TypeOf(dto)))
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
	apiErr := c.api.Take(spec)
	if apiErr != nil {
		return apiErr
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (c *echoController) PostSpec3(ctx echo.Context) error {
	var dto poolstep1.StepSpec
	bindErr := ctx.Bind(&dto)
	if bindErr != nil {
		c.log.Error("binding failed", slog.Any("dto", reflect.TypeOf(dto)))
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
