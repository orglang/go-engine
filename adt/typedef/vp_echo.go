package typedef

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	sdk "github.com/orglang/go-sdk/adt/descsem"

	"orglang/go-engine/lib/lf"
	"orglang/go-engine/lib/te"

	"orglang/go-engine/adt/descsem"
	"orglang/go-engine/adt/typeexp"
	"orglang/go-engine/adt/uniqsym"
)

// Adapter
type echoPresenter struct {
	api API
	ssr te.Renderer
	log *slog.Logger
}

func newEchoPresenter(a API, r te.Renderer, l *slog.Logger) *echoPresenter {
	name := slog.String("name", reflect.TypeFor[echoPresenter]().Name())
	return &echoPresenter{a, r, l.With(name)}
}

func cfgEchoPresenter(e *echo.Echo, p *echoPresenter) error {
	e.POST("/ssr/types", p.PostOne)
	e.GET("/ssr/types", p.GetMany)
	e.GET("/ssr/types/:id", p.GetOne)
	return nil
}

func (p *echoPresenter) PostOne(c echo.Context) error {
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
	qn, convertErr := uniqsym.ConvertFromString(dto.TypeQN)
	if convertErr != nil {
		p.log.Error("conversion failed", slog.Any("dto", dto))
		return convertErr
	}
	snap, createErr := p.api.Create(DefSpec{TypeQN: qn, TypeES: typeexp.OneSpec{}})
	if createErr != nil {
		return createErr
	}
	html, renderingErr := p.ssr.Render("view-one", ViewFromDefSnap(snap))
	if renderingErr != nil {
		p.log.Error("rendering failed", slog.Any("snap", snap))
		return renderingErr
	}
	p.log.Log(ctx, lf.LevelTrace, "posting succeed", slog.Any("ref", snap.DescRef))
	return c.HTMLBlob(http.StatusOK, html)
}

func (p *echoPresenter) GetMany(c echo.Context) error {
	refs, retrieveErr := p.api.RetreiveRefs()
	if retrieveErr != nil {
		return retrieveErr
	}
	html, renderingErr := p.ssr.Render("view-many", descsem.MsgFromRefs(refs))
	if renderingErr != nil {
		p.log.Error("rendering failed", slog.Any("refs", refs))
		return renderingErr
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func (p *echoPresenter) GetOne(c echo.Context) error {
	var dto sdk.SemRef
	bindErr := c.Bind(&dto)
	if bindErr != nil {
		p.log.Error("binding failed", slog.Any("dto", reflect.TypeOf(dto)))
		return bindErr
	}
	ctx := c.Request().Context()
	p.log.Log(ctx, lf.LevelTrace, "getting started", slog.Any("dto", dto))
	validateErr := dto.Validate()
	if validateErr != nil {
		p.log.Error("validation failed", slog.Any("dto", dto))
		return validateErr
	}
	ref, convertErr := descsem.MsgToRef(dto)
	if convertErr != nil {
		p.log.Error("conversion failed", slog.Any("dto", dto))
		return convertErr
	}
	snap, retrieveErr := p.api.RetrieveSnap(ref)
	if retrieveErr != nil {
		return retrieveErr
	}
	html, renderingErr := p.ssr.Render("view-one", ViewFromDefSnap(snap))
	if renderingErr != nil {
		p.log.Error("rendering failed", slog.Any("snap", snap))
		return renderingErr
	}
	p.log.Log(ctx, lf.LevelTrace, "getting succeed", slog.Any("ref", snap.DescRef))
	return c.HTMLBlob(http.StatusOK, html)
}
