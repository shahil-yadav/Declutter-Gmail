package web

import (
	"my-gmail-server/models"
	"my-gmail-server/pkg/app"
	"my-gmail-server/pkg/e"
	"my-gmail-server/templates"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	appG := app.Gin{C: c}
	users, err := models.ListUsers()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.ServeTempl(http.StatusOK, templates.IndexPage(users))
}
