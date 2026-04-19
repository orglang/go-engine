package termdec

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"

	sdk "github.com/orglang/go-sdk/adt/termsem"

	"orglang/go-engine/lib/lf"
	"orglang/go-engine/lib/te"

	"orglang/go-engine/adt/termsem"
	"orglang/go-engine/adt/uniqsym"
)

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
	e.POST("/ssr/decs", p.PostSpec)
	e.GET("/ssr/decs", p.GetRefs)
	e.GET("/ssr/decs/:id", p.GetSnap)
	return nil
}

func (p *echoPresenter) PostSpec(c echo.Context) error {
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
	ref, inceptionErr := p.api.Incept(qn)
	if inceptionErr != nil {
		return inceptionErr
	}
	html, renderingErr := p.ssr.Render("view-one", termsem.MsgFromRef(ref))
	if renderingErr != nil {
		p.log.Error("rendering failed", slog.Any("ref", ref))
		return renderingErr
	}
	p.log.Log(ctx, lf.LevelTrace, "posting succeed", slog.Any("ref", ref))
	return c.HTMLBlob(http.StatusOK, html)
}

func (p *echoPresenter) GetRefs(c echo.Context) error {
	refs, retrieveErr := p.api.RetreiveRefs()
	if retrieveErr != nil {
		return retrieveErr
	}
	html, renderingErr := p.ssr.Render("view-many", termsem.MsgFromRefs(refs))
	if renderingErr != nil {
		p.log.Error("rendering failed", slog.Any("refs", refs))
		return renderingErr
	}
	return c.HTMLBlob(http.StatusOK, html)
}

func (p *echoPresenter) GetSnap(c echo.Context) error {
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
	ref, convErr := termsem.MsgToRef(dto)
	if convErr != nil {
		p.log.Error("conversion failed", slog.Any("dto", dto))
		return convErr
	}
	snap, retrieveErr := p.api.RetrieveSnap(ref)
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
