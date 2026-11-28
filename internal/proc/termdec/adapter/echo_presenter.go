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
	termdec "orglang/go-engine/internal/proc/termdec/core"
)

type Presenter interface {
	PostSpec(*echo.Context) error
	GetRefs(*echo.Context) error
	GetSnap(*echo.Context) error
}

type echoPresenter struct {
	api termdec.API
	ssr te.Renderer
	log *slog.Logger
}

func newEchoPresenter(a termdec.API, r te.Renderer, l *slog.Logger) *echoPresenter {
	name := slog.String("name", reflect.TypeFor[echoPresenter]().Name())
	return &echoPresenter{a, r, l.With(name)}
}

func RegisterPresenterRoutes(e *echo.Echo, p Presenter) error {
	e.POST("/ssr/decs", p.PostSpec)
	e.GET("/ssr/decs", p.GetRefs)
	e.GET("/ssr/decs/:id", p.GetSnap)
	return nil
}
func cfgEchoPresenter(e *echo.Echo, p *echoPresenter) error { return RegisterPresenterRoutes(e, p) }

func (p *echoPresenter) PostSpec(c *echo.Context) error {
	var dto DecSpecVP
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
	qn, err := uniqsym.ConvertFromString(dto.TermQN)
	if err != nil {
		p.log.Error("conversion failed", slog.Any("dto", dto))
		return err
	}
	termID, progRef, err := p.api.Incept(ctx, qn)
	if err != nil {
		return err
	}
	view := DecSnapVP{}
	view.TermRef.TermID, view.TermRef.ProgID, view.TermRef.ProgRN = termID, progRef.ProgID, progRef.ProgRN
	html, err := p.ssr.Render("view-one", view)
	if err != nil {
		p.log.Error("rendering failed", slog.Any("ref", view))
		return err
	}
	p.log.Log(ctx, lf.LevelTrace, "posting succeed", slog.Any("termID", termID))
	return c.HTMLBlob(http.StatusOK, html)
}

func (p *echoPresenter) GetRefs(c *echo.Context) error {
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

func (p *echoPresenter) GetSnap(c *echo.Context) error {
	id, err := identity.ConvertFromString(c.Param("id"))
	if err != nil {
		p.log.Error("conversion failed", slog.String("id", c.Param("id")))
		return err
	}
	ctx := c.Request().Context()
	p.log.Log(ctx, lf.LevelTrace, "getting started", slog.Any("termID", id))
	snap, err := p.api.RetrieveSnap(ctx, id)
	if err != nil {
		return err
	}
	html, err := p.ssr.Render("view-one", ViewFromDecSnap(snap))
	if err != nil {
		p.log.Error("rendering failed", slog.Any("snap", snap))
		return err
	}
	p.log.Log(ctx, lf.LevelTrace, "getting succeed", slog.Any("termID", snap.TermID))
	return c.HTMLBlob(http.StatusOK, html)
}
