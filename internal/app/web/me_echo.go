package web

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v5"

	"orglang/go-engine/internal/lib/te"

	typedef "orglang/go-engine/internal/proc/typedef/core"
)

// Controller is the inbound HTTP port implemented by the Echo adapter.
type Controller interface {
	Home(*echo.Context) error
}

// RegisterRoutes registers the web application's HTTP routes.
// Controller construction is owned by the application container.
func RegisterRoutes(e *echo.Echo, h Controller) {
	e.GET("/", h.Home)
}

// Adapter
type controllerEcho struct {
	api typedef.API
	ssr te.Renderer
	log *slog.Logger
}

func newControllerEcho(a typedef.API, r te.Renderer, l *slog.Logger) *controllerEcho {
	name := slog.String("name", reflect.TypeFor[controllerEcho]().Name())
	return &controllerEcho{a, r, l.With(name)}
}

func cfgEchoController(e *echo.Echo, h *controllerEcho) {
	RegisterRoutes(e, h)
}

func (h *controllerEcho) Home(c *echo.Context) error {
	refs, err := h.api.RetreiveRefs()
	if err != nil {
		return err
	}
	html, err := h.ssr.Render("home.html", refs)
	if err != nil {
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}
