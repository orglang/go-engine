package adapter

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v5"

	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/te"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/uniqsym"
	typedef "orglang/go-engine/internal/proc/typedef/core"
)

type Presenter interface {
	PostOne(*echo.Context) error
	GetMany(*echo.Context) error
	GetOne(*echo.Context) error
}

type echoPresenter struct {
	api typedef.API
	ssr te.Renderer
	log *slog.Logger
}

func newEchoPresenter(a typedef.API, r te.Renderer, l *slog.Logger) *echoPresenter {
	name := slog.String("name", reflect.TypeFor[echoPresenter]().Name())
	return &echoPresenter{a, r, l.With(name)}
}

func RegisterPresenterRoutes(e *echo.Echo, p Presenter) error {
	e.POST("/ssr/types", p.PostOne)
	e.GET("/ssr/types", p.GetMany)
	e.GET("/ssr/types/:id", p.GetOne)
	return nil
}
func cfgEchoPresenter(e *echo.Echo, p *echoPresenter) error { return RegisterPresenterRoutes(e, p) }

func (p *echoPresenter) PostOne(c *echo.Context) error {
	var dto DefSpecVP
	if err := c.Bind(&dto); err != nil {
		p.log.Error("binding failed", slog.Any("dto", reflect.TypeOf(dto)))
		return err
	}
	ctx := c.Request().Context()
	p.log.Log(ctx, lf.LevelTrace, "posting started", slog.Any("dto", dto))
	if err := dto.Validate(); err != nil {
		p.log.Error("validation failed", slog.Any("dto", dto))
		return err
	}
	qn, err := uniqsym.ConvertFromString(dto.TypeQN)
	if err != nil {
		p.log.Error("conversion failed", slog.Any("dto", dto))
		return err
	}
	snap, err := p.api.Create(ctx, typedef.DefSpec{TypeQN: qn, TypeExprr: typedef.OneSpec{}})
	if err != nil {
		return err
	}
	html, err := p.ssr.Render("view-one", ViewFromDefSnap(snap))
	if err != nil {
		p.log.Error("rendering failed", slog.Any("snap", snap))
		return err
	}
	p.log.Log(ctx, lf.LevelTrace, "posting succeed", slog.Any("typeID", snap.TypeID))
	return c.HTMLBlob(http.StatusOK, html)
}

func (p *echoPresenter) GetMany(c *echo.Context) error {
	ids, err := p.api.RetreiveRefs(c.Request().Context())
	if err != nil {
		return err
	}
	html, err := p.ssr.Render("view-many", ids)
	if err != nil {
		p.log.Error("rendering failed", slog.Any("ids", ids))
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func (p *echoPresenter) GetOne(c *echo.Context) error {
	id, err := identity.ConvertFromString(c.Param("id"))
	if err != nil {
		p.log.Error("conversion failed", slog.String("id", c.Param("id")))
		return err
	}
	ctx := c.Request().Context()
	p.log.Log(ctx, lf.LevelTrace, "getting started", slog.Any("typeID", id))
	snap, err := p.api.RetrieveSnap(ctx, id)
	if err != nil {
		return err
	}
	html, err := p.ssr.Render("view-one", ViewFromDefSnap(snap))
	if err != nil {
		p.log.Error("rendering failed", slog.Any("snap", snap))
		return err
	}
	p.log.Log(ctx, lf.LevelTrace, "getting succeed", slog.Any("typeID", snap.TypeID))
	return c.HTMLBlob(http.StatusOK, html)
}
