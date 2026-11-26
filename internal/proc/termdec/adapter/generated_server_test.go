package adapter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"

	"orglang/go-engine/internal/adt/termsem"
	webapi "orglang/go-engine/internal/app/web/api"
)

type generatedDeclarationStub struct {
	webapi.ServerInterface
	controller *Controller
}

func (h generatedDeclarationStub) CreateProcessDeclaration(c *echo.Context) error {
	return c.NoContent(http.StatusCreated)
}
func (h generatedDeclarationStub) GetProcessDeclaration(c *echo.Context, id string, params webapi.GetProcessDeclarationParams) error {
	return c.NoContent(http.StatusOK)
}

func TestRegisterGeneratedRoutesRegistersOnlyDeclarationRoutes(t *testing.T) {
	e := echo.New()
	controller := &Controller{}
	webapi.RegisterHandlers(generatedRouteFilter{EchoRouter: e}, generatedDeclarationStub{controller: controller})

	for _, route := range []struct {
		method string
		path string
		want int
	}{
		{http.MethodPost, processDeclarationRoute, http.StatusCreated},
		{http.MethodGet, processDeclarationRoute+"/term-1", http.StatusOK},
		{http.MethodGet, "/api/v1/processes/term-1/definition", http.StatusNotFound},
		{http.MethodPost, "/api/v1/programs", http.StatusNotFound},
	} {
		req := httptest.NewRequestWithContext(context.Background(), route.method, route.path, nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		if rec.Code != route.want {
			t.Errorf("%s %s: got %d, want %d", route.method, route.path, rec.Code, route.want)
		}
	}
	_ = termsem.SemRef{}
}
