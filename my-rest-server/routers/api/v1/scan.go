package v1

import (
	"context"
	"fmt"
	"my-gmail-server/models"
	"my-gmail-server/pkg/app"
	"my-gmail-server/pkg/e"
	"my-gmail-server/services/jobs_service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mikestefanello/backlite"
)

func CreateScanJob(c *gin.Context) {
	appG := app.Gin{C: c}

	userIdStr := c.PostForm("user-id")
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.STATUS_BAD_REQUEST, err.Error())
		return
	}

	user, err := models.ListUserById(userId)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	if err := jobs_service.CreateScanJob(userId); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS,
		fmt.Sprintf("Started the job to scan %v's mailbox (%v)", user.FullName, user.Email),
	)
}

type JobStatus struct {
	JobId     string
	IsPending bool
	IsSuccess bool
	IsError   bool
}

type ScanJobResults struct {
	JobStatus
	Results []models.MailFolder
}

func ListActiveScanJobInfo(c *gin.Context) {
	appG := app.Gin{C: c}
	param := c.Param("id")

	userId, err := strconv.Atoi(param)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.STATUS_BAD_REQUEST, err.Error())
		return
	}

	job, err := models.GetActiveScanJob(userId)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	status, err := jobs_service.Client.Status(context.Background(), job)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	var scanJobResults ScanJobResults

	scanJobResults.JobId = job

	switch status {
	case backlite.TaskStatusRunning:
		scanJobResults.IsPending = true

	case backlite.TaskStatusSuccess:
		scanJobResults.IsSuccess = true

	case backlite.TaskStatusFailure:
		scanJobResults.IsError = true
	}

	if scanJobResults.IsSuccess {
		scanJobResults.Results, err = models.ListMailFolders(userId)
		if err != nil {
			appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
			return
		}

	}

	appG.Response(http.StatusOK, e.SUCCESS, scanJobResults)
}
