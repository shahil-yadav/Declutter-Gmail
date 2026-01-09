package app

import (
	"my-gmail-server/pkg/e"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	g.C.JSON(httpCode, Response{
		Code: errCode,
		Msg:  e.GetMsg(errCode),
		Data: data,
	})
}

// Serve Templ
func (g *Gin) ServeTempl(httpCode int, component templ.Component) {
	g.C.HTML(httpCode, "", component)
}
