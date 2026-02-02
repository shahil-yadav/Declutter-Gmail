// this file is responsible for creating scan jobs

package jobs_service

import (
	"context"
	"errors"
	"my-gmail-server/models"
	"my-gmail-server/pkg/utils"
	"my-gmail-server/services/auth_service"
	"time"

	"github.com/mikestefanello/backlite"
)

type ScanJob struct {
	// p.k of `users` table
	UserId int
}

func (job ScanJob) Config() backlite.QueueConfig {
	config := backlite.QueueConfig{
		Name:        "scan-job-queue",
		MaxAttempts: 1,
		Backoff:     5 * time.Second,
		Retention:   &backlite.Retention{Data: &backlite.RetainData{}},
	}

	return config
}

func ScanJobWork(c context.Context, job ScanJob) error {
	token, err := models.GetUserTokens(job.UserId)
	if err != nil {
		return err
	}

	gService, err := auth_service.CreateGmailServiceWithTokens(token)
	if err != nil {
		return err
	}

	profile, err := gService.Users.GetProfile("me").Do()
	if err != nil {
		return err
	}

	// Delete every mails related user-id, each scan job start a new fresh sweep of mailbox
	if err := models.DeleteMailsByUserId(job.UserId); err != nil {
		return err
	}

	gMessages := utils.GetAllMessageIds(gService)

	threads := 20
	mails := utils.DistributedScan(gService, profile, job.UserId, gMessages, threads)

	/* Synchronus computing earlier */
	// var mails []models.Mail
	// for _, message := range gMessages {
	// 	message, err := gService.Users.Messages.Get("me", message.Id).Format("metadata").Do()
	// 	// todo: improve this logic
	// 	if err != nil {
	// 		fmt.Println("Failed to fetch message:", message.Id)
	// 		continue
	// 	}

	// 	senderEmail, err := utils.ExtractSenderAdressFromGmailMessage(message)
	// 	if err != nil {
	// 		fmt.Println("Failed to extract sender email from message", message.Id)
	// 		continue
	// 	}

	// 	date, err := utils.ExtractDateFromGmailMessage(message)
	// 	if err != nil {
	// 		fmt.Println("Failed to extract date in UTC from message", message.Id)
	// 		//
	// 		continue
	// 	}

	// 	mails = append(mails, models.Mail{
	// 		Id:           message.Id,
	// 		AccountEmail: profile.EmailAddress,
	// 		SenderEmail:  senderEmail,
	// 		Snippet:      message.Snippet,
	// 		Date:         date,
	// 		UserId:       job.UserId,
	// 	})
	// }

	err = models.AddMails(mails)
	if err != nil {
		return err
	}

	return nil
}

func NewScanJob(userId int) ScanJob {
	return ScanJob{userId}
}

func CreateScanJob(userId int) error {
	// if err := models.DeactivateScanJobs(userId); err != nil {
	// 	return err
	// }

	job := NewScanJob(userId)
	jobIds, err := Client.Add(job).Save()
	if err != nil {
		return err
	}
	if len(jobIds) > 1 {
		return errors.New("i fcked with client job queueing")
	}

	jobId := jobIds[0]
	if err := models.AddScanJob(userId, jobId); err != nil {
		return err
	}

	return nil
}
