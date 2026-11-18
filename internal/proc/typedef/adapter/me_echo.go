package adapter

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	sdktypedef "github.com/orglang/go-sdk/proc/typedef"

	"orglang/go-engine/internal/lib/lf"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/typesem"
	typedef "orglang/go-engine/internal/proc/typedef/core"
)

// Controller is the inbound HTTP port implemented by the Echo adapter.
type Controller interface {
	PostSpec(echo.Context) error
	GetSnap(echo.Context) error
	PatchOne(echo.Context) error
}

// Server-side primary adapter
type typeDefRefDTO struct {
	TypeID string `json:"type_id"`
	TypeRN int64  `json:"type_rn"`
}

type controllerEcho struct {
	api typedef.API
	log *slog.Logger
}

func newControllerEcho(a typedef.API, l *slog.Logger) *controllerEcho {
	name := slog.String("name", reflect.TypeFor[controllerEcho]().Name())
	return &controllerEcho{a, l.With(name)}
}

// RegisterRoutes attaches the type-definition API routes to an Echo server.
func RegisterRoutes(e *echo.Echo, h Controller) error {
	e.POST("/api/v1/types", h.PostSpec)
	e.GET("/api/v1/types/:id", h.GetSnap)
	e.PATCH("/api/v1/types/:id", h.PatchOne)
	return nil
}

func cfgEchoController(e *echo.Echo, h *controllerEcho) error {
	return RegisterRoutes(e, h)
}

func (h *controllerEcho) PostSpec(c echo.Context) error {
	var dto sdktypedef.DefSpec
	bindErr := c.Bind(&dto)
	if bindErr != nil {
		h.log.Error("binding failed", slog.Any("dto", reflect.TypeOf(dto)))
		return bindErr
	}
	ctx := c.Request().Context()
	h.log.Log(ctx, lf.LevelTrace, "posting started", slog.Any("dto", dto))
	validateErr := dto.Validate()
	if validateErr != nil {
		h.log.Error("validation failed", slog.Any("dto", dto))
		return validateErr
	}
	spec, convErr := MsgToDefSpec(dto)
	if convErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return convErr
	}
	snap, createErr := h.api.Create(spec)
	if createErr != nil {
		return createErr
	}
	h.log.Log(ctx, lf.LevelTrace, "posting succeed", slog.Any("ref", snap.TypeRef))
	return c.JSON(http.StatusCreated, MsgFromDefSnap(snap))
}

func (h *controllerEcho) GetSnap(c echo.Context) error {
	var dto typeDefRefDTO
	bindErr := c.Bind(&dto)
	if bindErr != nil {
		h.log.Error("binding failed", slog.Any("dto", reflect.TypeOf(dto)))
		return bindErr
	}
	typeID, convErr := identity.ConvertFromString(dto.TypeID)
	if convErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return echo.NewHTTPError(http.StatusBadRequest, convErr.Error())
	}
	if dto.TypeRN < 1 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid revision")
	}
	ref := typesem.SemRef{TypeID: typeID, TypeRN: seqnum.ConvertFromInt(dto.TypeRN)}
	snap, retrieveErr := h.api.RetrieveSnap(ref)
	if retrieveErr != nil {
		return retrieveErr
	}
	return c.JSON(http.StatusOK, MsgFromDefSnap(snap))
}

func (h *controllerEcho) PatchOne(c echo.Context) error {
	var dto sdktypedef.DefSnap
	bindErr := c.Bind(&dto)
	if bindErr != nil {
		h.log.Error("binding failed", slog.Any("dto", reflect.TypeOf(dto)))
		return bindErr
	}
	ctx := c.Request().Context()
	h.log.Log(ctx, lf.LevelTrace, "patching started", slog.Any("dto", dto))
	validateErr := dto.Validate()
	if validateErr != nil {
		h.log.Error("validation failed", slog.Any("dto", dto))
		return validateErr
	}
	reqSnap, convErr := MsgToDefSnap(dto)
	if convErr != nil {
		h.log.Error("conversion failed", slog.Any("dto", dto))
		return convErr
	}
	resSnap, modificationErr := h.api.Modify(reqSnap)
	if modificationErr != nil {
		return modificationErr
	}
	h.log.Log(ctx, lf.LevelTrace, "patching succeed", slog.Any("ref", resSnap.TypeRef))
	return c.JSON(http.StatusOK, MsgFromDefSnap(resSnap))
}
