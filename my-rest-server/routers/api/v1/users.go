package v1

import (
	"my-gmail-server/models"
	"my-gmail-server/pkg/app"
	"my-gmail-server/pkg/e"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListUsers(c *gin.Context) {
	appG := app.Gin{C: c}

	users, err := models.ListUsers()
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, users)
}

func ListUserById(c *gin.Context) {
	appG := app.Gin{C: c}
	param := c.Param("id")

	userId, err := strconv.Atoi(param)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.STATUS_BAD_REQUEST, err.Error())
		return
	}

	users, err := models.ListUserById(userId)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, users)
}
