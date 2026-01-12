package v1

import (
	"my-gmail-server/pkg/app"
	"my-gmail-server/pkg/e"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckHealth(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, e.SUCCESS, "perfectly healthy")
}
