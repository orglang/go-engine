package adapter

import (
	"log/slog"
	"net/http"
	"reflect"
	"strconv"

	"github.com/labstack/echo/v5"

	"orglang/go-engine/internal/lib/lf"
	"orglang/go-engine/internal/lib/te"

	"orglang/go-engine/internal/adt/identity"
	"orglang/go-engine/internal/adt/seqnum"
	"orglang/go-engine/internal/adt/termsem"
	"orglang/go-engine/internal/adt/uniqsym"

	termdec "orglang/go-engine/internal/proc/termdec/core"
)

// Presenter is the server-rendered HTTP port implemented by the Echo adapter.
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

// RegisterPresenterRoutes attaches the term-declaration SSR routes to Echo.
func RegisterPresenterRoutes(e *echo.Echo, p Presenter) error {
	e.POST("/ssr/decs", p.PostSpec)
	e.GET("/ssr/decs", p.GetRefs)
	e.GET("/ssr/decs/:id", p.GetSnap)
	return nil
}

func cfgEchoPresenter(e *echo.Echo, p *echoPresenter) error {
	return RegisterPresenterRoutes(e, p)
}

func (p *echoPresenter) PostSpec(c *echo.Context) error {
	var dto DecSpecVP
	bindErr := c.Bind(&dto)
	if bindErr != nil {
		p.log.Error("binding failed", slog.Any("dto", reflect.TypeOf(dto)))
		return bindErr
	}
	ctx := c.Request().Context()
	p.log.Log(ctx, lf.LevelTrace, "posting started", slog.Any("dto", dto))
	validateErr := dto.Validate()
	if validateErr != nil {
		p.log.Error("validation failed", slog.Any("dto", dto))
		return validateErr
	}
	qn, convErr := uniqsym.ConvertFromString(dto.TermQN)
	if convErr != nil {
		p.log.Error("conversion failed", slog.Any("dto", dto))
		return convErr
	}
	ref, inceptionErr := p.api.Incept(ctx, qn)
	if inceptionErr != nil {
		return inceptionErr
	}
	html, renderingErr := p.ssr.Render("view-one", ref)
	if renderingErr != nil {
		p.log.Error("rendering failed", slog.Any("ref", ref))
		return renderingErr
	}
	p.log.Log(ctx, lf.LevelTrace, "posting succeed", slog.Any("ref", ref))
	return c.HTMLBlob(http.StatusOK, html)
}

func (p *echoPresenter) GetRefs(c *echo.Context) error {
	refs, retrieveErr := p.api.RetreiveRefs(c.Request().Context())
	if retrieveErr != nil {
		return retrieveErr
	}
	html, renderingErr := p.ssr.Render("view-many", refs)
	if renderingErr != nil {
		p.log.Error("rendering failed", slog.Any("refs", refs))
		return renderingErr
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func (p *echoPresenter) GetSnap(c *echo.Context) error {
	id, convErr := identity.ConvertFromString(c.Param("id"))
	if convErr != nil {
		p.log.Error("conversion failed", slog.String("id", c.Param("id")))
		return convErr
	}
	var rn int64
	if value := c.QueryParam("rn"); value != "" {
		rn, convErr = strconv.ParseInt(value, 10, 64)
		if convErr != nil {
			p.log.Error("invalid revision", slog.String("rn", value))
			return echo.NewHTTPError(http.StatusBadRequest, "invalid revision").Wrap(convErr)
		}
	}
	ref := termsem.SemRef{TermID: id, DescRN: seqnum.ConvertFromInt(rn)}
	ctx := c.Request().Context()
	p.log.Log(ctx, lf.LevelTrace, "getting started", slog.Any("ref", ref))
	snap, retrieveErr := p.api.RetrieveSnap(ctx, ref)
	if retrieveErr != nil {
		return retrieveErr
	}
	html, renderingErr := p.ssr.Render("view-one", ViewFromDecSnap(snap))
	if renderingErr != nil {
		p.log.Error("rendering failed", slog.Any("snap", snap))
		return renderingErr
	}
	p.log.Log(ctx, lf.LevelTrace, "getting succeed", slog.Any("decRef", snap.TermRef))
	return c.HTMLBlob(http.StatusOK, html)
}
