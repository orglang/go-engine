package poolstep

import (
	"log/slog"
	"net/http"
	"orglang/go-engine/adt/implsem"
	"reflect"

	"github.com/labstack/echo/v4"

	"github.com/orglang/go-sdk/adt/poolstep"
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

func cfgEchoController(e *echo.Echo, h *echoController) error {
	e.POST("/api/v1/pools/execs/steps", h.PostSpec)
	e.POST("/api/v1/pools/execs/spawns", h.PostSpec2)
	return nil
}

func (h *echoController) PostSpec(c echo.Context) error {
	var dto poolstep.StepSpec
	bindErr := c.Bind(&dto)
	if bindErr != nil {
		h.log.Error("binding failed", slog.Any("dto", reflect.TypeOf(dto)))
		return bindErr
	}
	validateErr := dto.Validate()
	if validateErr != nil {
		h.log.Error("validation failed", slog.Any("dto", dto))
		return validateErr
	}
	spec, convertErr := MsgToStepSpec(dto)
	if convertErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return convertErr
	}
	takeErr := h.api.Take(spec)
	if takeErr != nil {
		return takeErr
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *echoController) PostSpec2(c echo.Context) error {
	var dto poolstep.StepSpec
	bindErr := c.Bind(&dto)
	if bindErr != nil {
		h.log.Error("binding failed", slog.Any("dto", reflect.TypeOf(dto)))
		return bindErr
	}
	validateErr := dto.Validate()
	if validateErr != nil {
		h.log.Error("validation failed", slog.Any("dto", dto))
		return validateErr
	}
	spec, convertErr := MsgToStepSpec(dto)
	if convertErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return convertErr
	}
	ref, takeErr := h.api.Spawn(spec)
	if takeErr != nil {
		return takeErr
	}
	return c.JSON(http.StatusCreated, implsem.MsgFromRef(ref))
}
