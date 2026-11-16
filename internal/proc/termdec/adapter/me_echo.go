package adapter

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	sdk "github.com/orglang/go-sdk/adt/termsem"
	sdktermdec "github.com/orglang/go-sdk/proc/termdec"

	"orglang/go-engine/internal/adt/termsem"
	webapi "orglang/go-engine/internal/app/web/api"
	termdec "orglang/go-engine/internal/proc/termdec/core"
)

// Controller is the inbound HTTP port implemented by the Echo adapter.
type Controller interface {
	PostSpec(echo.Context) error
	GetSnap(echo.Context) error
	CreateProcessDeclaration(echo.Context) error
}

// Server-side primary adapter
type controllerEcho struct {
	// Only operations explicitly mounted by the generated route filter are
	// dispatched; embedding keeps the generated interface forward-compatible.
	webapi.ServerInterface
	api termdec.API
	log *slog.Logger
}

var _ webapi.ServerInterface = (*controllerEcho)(nil)

func newControllerEcho(a termdec.API, l *slog.Logger) *controllerEcho {
	name := slog.String("name", reflect.TypeFor[controllerEcho]().Name())
	return &controllerEcho{api: a, log: l.With(name)}
}

// RegisterRoutes attaches the term-declaration API routes to an Echo server.
func RegisterRoutes(e *echo.Echo, h Controller) error {
	e.POST("/api/v1/procs/decs", h.PostSpec)
	e.GET("/api/v1/procs/decs/:id", h.GetSnap)
	return nil
}

func cfgEchoController(e *echo.Echo, h *controllerEcho) error {
	return RegisterRoutes(e, h)
}

func (h *controllerEcho) PostSpec(c echo.Context) error {
	var dto sdktermdec.DecSpec
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
	spec, convErr := MsgToDecSpec(dto)
	if convErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return convErr
	}
	snap, createErr := h.api.Create(spec)
	if createErr != nil {
		return createErr
	}
	return c.JSON(http.StatusCreated, MsgFromDecSnap(snap))
}

func (h *controllerEcho) GetSnap(c echo.Context) error {
	var dto sdk.SemRef
	bindErr := c.Bind(&dto)
	if bindErr != nil {
		h.log.Error("binding failed", slog.Any("dto", reflect.TypeOf(dto)))
		return bindErr
	}
	ref, convErr := termsem.MsgToRef(dto)
	if convErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return convErr
	}
	snap, retrieveErr := h.api.RetrieveSnap(ref)
	if retrieveErr != nil {
		return retrieveErr
	}
	return c.JSON(http.StatusOK, MsgFromDecSnap(snap))
}

func (h *controllerEcho) CreateProcessDeclaration(c echo.Context) error {
	var dto webapi.ProcDecSpec
	if err := c.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	spec, err := DecSpecFromGenerated(dto)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	snap, err := h.api.Create(spec)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, DecSnapToGenerated(snap))
}
