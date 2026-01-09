package gintemplrenderer

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"

	"github.com/a-h/templ"
)

// setup HTML rendering
// im going with go+ssr, no spa yet
func Setup(e *gin.Engine) {
	e.LoadHTMLGlob("templates/*")
	ginHtmlRenderer := e.HTMLRender
	e.HTMLRender = &HTMLTemplRenderer{FallbackHtmlRenderer: ginHtmlRenderer}
}

type HTMLTemplRenderer struct {
	FallbackHtmlRenderer render.HTMLRender
}

func (r *HTMLTemplRenderer) Instance(s string, d any) render.Render {
	templData, ok := d.(templ.Component)
	if !ok {
		if r.FallbackHtmlRenderer != nil {
			return r.FallbackHtmlRenderer.Instance(s, d)
		}
	}
	return &Renderer{
		Ctx:       context.Background(),
		Status:    -1,
		Component: templData,
	}
}

type Renderer struct {
	Ctx       context.Context
	Status    int
	Component templ.Component
}

func (t Renderer) Render(w http.ResponseWriter) error {
	t.WriteContentType(w)
	if t.Status != -1 {
		w.WriteHeader(t.Status)
	}
	if t.Component != nil {
		return t.Component.Render(t.Ctx, w)
	}
	return nil
}

func (t Renderer) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}
