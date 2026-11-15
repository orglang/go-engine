package web

import (
	"github.com/labstack/echo/v4"

	webapi "orglang/go-engine/internal/app/web/api"
)

// cfgGeneratedServer registers the contract-generated Echo routes on the
// application's existing Echo instance. The generated API is not a second
// listener or server.
func cfgGeneratedServer(e *echo.Echo, h webapi.ServerInterface) {
	webapi.RegisterHandlers(e, h)
}
