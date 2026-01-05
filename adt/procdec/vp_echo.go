package procdec

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	"orglang/orglang/adt/identity"
	"orglang/orglang/adt/qualsym"

	"orglang/orglang/lib/lf"
	"orglang/orglang/lib/te"
)

// Adapter
type presenterEcho struct {
	api API
	ssr te.Renderer
	log *slog.Logger
}

func newPresenterEcho(a API, r te.Renderer, l *slog.Logger) *presenterEcho {
	name := slog.String("name", reflect.TypeFor[presenterEcho]().Name())
	return &presenterEcho{a, r, l.With(name)}
}

func cfgPresenterEcho(e *echo.Echo, p *presenterEcho) error {
	e.POST("/ssr/declarations", p.PostOne)
	e.GET("/ssr/declarations", p.GetMany)
	e.GET("/ssr/declarations/:id", p.GetOne)
	return nil
}

func (p *presenterEcho) PostOne(c echo.Context) error {
	var dto DecSpecVP
	err := c.Bind(&dto)
	if err != nil {
		p.log.Error("dto binding failed")
		return err
	}
	ctx := c.Request().Context()
	p.log.Log(ctx, lf.LevelTrace, "root posting started", slog.Any("dto", dto))
	err = dto.Validate()
	if err != nil {
		p.log.Error("dto validation failed")
		return err
	}
	ns, err := qualsym.ConvertFromString(dto.ProcNS)
	if err != nil {
		p.log.Error("dto parsing failed")
		return err
	}
	ref, err := p.api.Incept(ns.New(dto.ProcSN))
	if err != nil {
		p.log.Error("root creation failed")
		return err
	}
	html, err := p.ssr.Render("view-one", ViewFromDecRef(ref))
	if err != nil {
		p.log.Error("view rendering failed")
		return err
	}
	p.log.Log(ctx, lf.LevelTrace, "root posting succeed", slog.Any("ref", ref))
	return c.HTMLBlob(http.StatusOK, html)
}

func (p *presenterEcho) GetMany(c echo.Context) error {
	refs, err := p.api.RetreiveRefs()
	if err != nil {
		p.log.Error("refs retrieval failed")
		return err
	}
	html, err := p.ssr.Render("view-many", ViewFromDecRefs(refs))
	if err != nil {
		p.log.Error("view rendering failed")
		return err
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func (p *presenterEcho) GetOne(c echo.Context) error {
	var dto IdentME
	err := c.Bind(&dto)
	if err != nil {
		p.log.Error("dto binding failed")
		return err
	}
	ctx := c.Request().Context()
	p.log.Log(ctx, lf.LevelTrace, "root getting started", slog.Any("dto", dto))
	err = dto.Validate()
	if err != nil {
		p.log.Error("dto validation failed")
		return err
	}
	id, err := identity.ConvertFromString(dto.DecID)
	if err != nil {
		p.log.Error("dto mapping failed")
		return err
	}
	snap, err := p.api.RetrieveSnap(id)
	if err != nil {
		p.log.Error("snap retrieval failed")
		return err
	}
	html, err := p.ssr.Render("view-one", ViewFromDecSnap(snap))
	if err != nil {
		p.log.Error("view rendering failed")
		return err
	}
	p.log.Log(ctx, lf.LevelTrace, "root getting succeed", slog.Any("id", snap.DecID))
	return c.HTMLBlob(http.StatusOK, html)
}
