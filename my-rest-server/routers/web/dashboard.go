package web

import (
	"context"
	"my-gmail-server/models"
	"my-gmail-server/pkg/app"
	"my-gmail-server/pkg/e"
	"my-gmail-server/services/jobs_service"
	"my-gmail-server/templates"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mikestefanello/backlite"
)

func Dashboard(c *gin.Context) {
	appG := app.Gin{C: c}

	param := c.Param("user")
	userId, err := strconv.Atoi(param)
	if err != nil {
		appG.Response(http.StatusUnauthorized, e.STATUS_UNAUTHORIZED, userId)
		return
	}

	activeScanJob, err := models.GetActiveScanJob(userId)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	user, err := models.ListUserById(userId)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	status, err := jobs_service.Client.Status(context.Background(), activeScanJob)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	var s string
	switch status {
	case backlite.TaskStatusRunning:
		s = models.Pending

	case backlite.TaskStatusFailure:
		s = models.Failed
	}

	mailFolders, err := models.ListMailFolders(userId)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.ServeTempl(http.StatusOK, templates.DashboardPage(user.FullName, s, userId, mailFolders))
}
