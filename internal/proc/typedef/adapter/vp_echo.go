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
	"orglang/go-engine/internal/adt/typesem"
	"orglang/go-engine/internal/adt/uniqsym"
	typedef "orglang/go-engine/internal/proc/typedef/core"
)

// Presenter is the server-rendered HTTP port implemented by the Echo adapter.
type Presenter interface {
	PostOne(*echo.Context) error
	GetMany(*echo.Context) error
	GetOne(*echo.Context) error
}

// Adapter
type echoPresenter struct {
	api typedef.API
	ssr te.Renderer
	log *slog.Logger
}

func newEchoPresenter(a typedef.API, r te.Renderer, l *slog.Logger) *echoPresenter {
	name := slog.String("name", reflect.TypeFor[echoPresenter]().Name())
	return &echoPresenter{a, r, l.With(name)}
}

// RegisterPresenterRoutes attaches the type-definition SSR routes to Echo.
func RegisterPresenterRoutes(e *echo.Echo, p Presenter) error {
	e.POST("/ssr/types", p.PostOne)
	e.GET("/ssr/types", p.GetMany)
	e.GET("/ssr/types/:id", p.GetOne)
	return nil
}

func cfgEchoPresenter(e *echo.Echo, p *echoPresenter) error {
	return RegisterPresenterRoutes(e, p)
}

func (p *echoPresenter) PostOne(c *echo.Context) error {
	var dto DefSpecVP
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
	qn, convErr := uniqsym.ConvertFromString(dto.TypeQN)
	if convErr != nil {
		p.log.Error("conversion failed", slog.Any("dto", dto))
		return convErr
	}
	snap, createErr := p.api.Create(ctx, typedef.DefSpec{TypeQN: qn, TypeExp: typedef.OneSpec{}})
	if createErr != nil {
		return createErr
	}
	html, renderingErr := p.ssr.Render("view-one", ViewFromDefSnap(snap))
	if renderingErr != nil {
		p.log.Error("rendering failed", slog.Any("snap", snap))
		return renderingErr
	}
	p.log.Log(ctx, lf.LevelTrace, "posting succeed", slog.Any("ref", snap.TypeRef))
	return c.HTMLBlob(http.StatusOK, html)
}

func (p *echoPresenter) GetMany(c *echo.Context) error {
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

func (p *echoPresenter) GetOne(c *echo.Context) error {
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
	ref := typesem.SemRef{TypeID: id, DescRN: seqnum.ConvertFromInt(rn)}
	ctx := c.Request().Context()
	p.log.Log(ctx, lf.LevelTrace, "getting started", slog.Any("ref", ref))
	snap, retrieveErr := p.api.RetrieveSnap(ctx, ref)
	if retrieveErr != nil {
		return retrieveErr
	}
	html, renderingErr := p.ssr.Render("view-one", ViewFromDefSnap(snap))
	if renderingErr != nil {
		p.log.Error("rendering failed", slog.Any("snap", snap))
		return renderingErr
	}
	p.log.Log(ctx, lf.LevelTrace, "getting succeed", slog.Any("ref", snap.TypeRef))
	return c.HTMLBlob(http.StatusOK, html)
}
