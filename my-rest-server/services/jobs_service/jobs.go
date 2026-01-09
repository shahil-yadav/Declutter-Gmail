package jobs_service

import (
	"context"

	"my-gmail-server/models"

	"github.com/mikestefanello/backlite"
	"golang.org/x/oauth2"
)

type DownloadEmailsJob struct {
	AccountEmail string

	// had to use oauth2.Token since this is serializable
	Tokens oauth2.Token
}

type JobStatus map[string][]string

func CollectScanJobsWithStatus(userId int) (JobStatus, error) {
	jobStatus := make(JobStatus)

	jobIds, err := models.ListScanJobs(userId)
	if err != nil {
		return nil, err
	}

	var pending, success, failed []string
	for _, jobId := range jobIds {
		jobStatus, err := Client.Status(context.Background(), jobId)
		if err != nil {
			return nil, err
		}

		switch jobStatus {
		case backlite.TaskStatusPending:
			pending = append(pending, jobId)

		case backlite.TaskStatusSuccess:
			success = append(pending, jobId)

		case backlite.TaskStatusFailure:
			failed = append(pending, jobId)

		}
	}

	jobStatus[models.Pending] = pending
	jobStatus[models.Success] = success
	jobStatus[models.Failed] = failed

	return jobStatus, nil
}
