package adapter

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v5"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	webapi "orglang/go-engine/internal/app/web/api"
	termdef "orglang/go-engine/internal/pool/termdef/core"
)

// GeneratedController is the inbound HTTP port for pool term definitions.
type GeneratedController interface {
	CreatePoolDeclaration(*echo.Context) error
}

type controllerEcho struct {
	api termdef.API
	log *slog.Logger
}

var _ GeneratedController = (*controllerEcho)(nil)

func newControllerEcho(a termdef.API, l *slog.Logger) *controllerEcho {
	name := slog.String("name", reflect.TypeFor[controllerEcho]().Name())
	return &controllerEcho{api: a, log: l.With(name)}
}

func (h *controllerEcho) CreatePoolDeclaration(c *echo.Context) error {
	var dto webapi.PoolDecSpec
	if err := c.Bind(&dto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	spec, err := DecSpecFromGenerated(dto)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	ref, err := h.api.Create(spec)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, webapi.AnyTermRef{
		TermId: identity.ConvertToString(ref.TermID),
		TermRn: seqnum.ConvertToInt(ref.TermRN),
	})
}
