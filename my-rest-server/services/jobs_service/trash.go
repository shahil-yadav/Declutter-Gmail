package jobs_service

import (
	"context"
	"my-gmail-server/models"
	"my-gmail-server/pkg/utils"
	"time"

	"github.com/mikestefanello/backlite"
)

type TrashJob struct {
	// p.k of `users` table
	UserId int

	// senders addresses to delete
	Senders []string
}

func (job TrashJob) Config() backlite.QueueConfig {
	return backlite.QueueConfig{
		Name:        "trash-job-queue",
		MaxAttempts: 1,
		Backoff:     5 * time.Second,
		Retention:   &backlite.Retention{Data: &backlite.RetainData{}},
	}
}

func TrashJobWork(c context.Context, job TrashJob) error {
	time.Sleep(1 * time.Minute)
	return nil

	gService, err := utils.CreateGmailService(job.UserId)
	if err != nil {
		return err
	}

	var toTrashIds []string
	for _, sender := range job.Senders {
		mails, err := models.ListMailsBySender(sender, job.UserId)
		if err != nil {
			return err
		}

		for _, mail := range mails {
			toTrashIds = append(toTrashIds, mail.Id)
		}
	}

	for _, msgId := range toTrashIds {
		if err := utils.TrashGmailMessage(gService, msgId); err != nil {
			return err
		}

		if err := models.DeleteMail(msgId); err != nil {
			// im going to silently ignore the error of untrashing
			utils.UntrashGmailMessage(gService, msgId)
			return err
		}
	}

	return nil
}

func CreateTrashJob(userId int, senders []string) TrashJob {
	return TrashJob{UserId: userId, Senders: senders}
}
