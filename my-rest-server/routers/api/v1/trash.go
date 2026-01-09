package v1

import (
	"context"
	"database/sql"
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

type trashForm struct {
	Senders []string `form:"sender[]"`
	UserId  int      `form:"user-id"`
}

type TrashJobResults struct {
	JobStatus
	NoExistingJobs bool // user is new to the platform and hasn't create any trash job
}

func ListActiveTrashJobInfo(c *gin.Context) {
	var trashJobResults TrashJobResults

	appG := app.Gin{C: c}
	param := c.Param("id")

	userId, err := strconv.Atoi(param)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.STATUS_BAD_REQUEST, err.Error())
		return
	}

	job, err := models.GetActiveTrashJob(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			trashJobResults.NoExistingJobs = true
		} else {
			appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
			return
		}
	}

	status, err := jobs_service.Client.Status(context.Background(), job)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return

	}

	trashJobResults.JobId = job

	switch status {
	case backlite.TaskStatusRunning:
		trashJobResults.IsPending = true

	case backlite.TaskStatusSuccess:
		trashJobResults.IsSuccess = true

	case backlite.TaskStatusFailure:
		trashJobResults.IsError = true
	}

	appG.Response(http.StatusOK, e.SUCCESS, trashJobResults)
}

func CreateTrashJob(c *gin.Context) {
	appG := app.Gin{C: c}

	var trashForm trashForm
	err := c.Bind(&trashForm)
	if err != nil {
		appG.Response(http.StatusBadRequest, e.STATUS_BAD_REQUEST, err.Error())
		return
	}

	job := jobs_service.CreateTrashJob(trashForm.UserId, trashForm.Senders)
	jobIds, err := jobs_service.Client.Add(job).Save()
	if err != nil || len(jobIds) > 1 {
		appG.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	jobId := jobIds[0]
	if err := models.AddTrashJob(jobId, trashForm.UserId); err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR, err.Error())
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS,
		fmt.Sprintf("Started the job to trash senders (%v) of user-id:%d", job.Senders, job.UserId),
	)
}
