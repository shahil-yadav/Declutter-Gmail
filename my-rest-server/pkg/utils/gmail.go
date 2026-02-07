package utils

import (
	"context"
	"fmt"
	"my-gmail-server/models"
	"my-gmail-server/services/auth_service"
	"net/mail"
	"time"

	"google.golang.org/api/gmail/v1"
)

func CreateGmailService(userId int) (*gmail.Service, error) {
	token, err := models.GetUserTokens(userId)
	if err != nil {
		return nil, err
	}

	gService, err := auth_service.CreateGmailServiceWithTokens(token)
	if err != nil {
		return nil, err
	}

	return gService, nil
}

func GetMessageIDs(gService *gmail.Service, count int64) ([]*gmail.Message, error) {
	if count > 500 {
		panic("GetMessageIDs is not designed to handle the count with more than 500")
	}

	lmr, err := gService.Users.Messages.List("me").MaxResults(count).Do()
	if err != nil {
		return nil, err
	}

	return lmr.Messages, nil
}

func GetAllMessageIds(gService *gmail.Service) []*gmail.Message {
	var (
		messages []*gmail.Message
		ctx      = context.Background()
	)

	gService.Users.Messages.List("me").MaxResults(500).Pages(
		ctx,
		func(lmr *gmail.ListMessagesResponse) error {
			messages = append(messages, lmr.Messages...)
			return nil
		},
	)

	return messages
}

func ExtractDateFromGmailMessage(message *gmail.Message) (time.Time, error) {
	var t time.Time
	var err error

	for _, header := range message.Payload.Headers {
		if header.Name == "Date" {
			t, err = mail.ParseDate(header.Value)
			t = t.UTC()

			break
		}
	}
	return t, err
}

func ExtractSenderAdressFromGmailMessage(message *gmail.Message) (string, error) {
	var fromHeader string
	for _, header := range message.Payload.Headers {
		if header.Name == "From" {
			fromHeader = header.Value
			// no need to go further in search of headers, i got what i needed
			break
		}
	}

	// google's official library for mail utils
	parsed, err := mail.ParseAddress(fromHeader)
	if err != nil {
		return "", err
	}

	// Works 99.99% in most cases hehe,
	return parsed.Address, nil
}

func TrashGmailMessage(gService *gmail.Service, id string) error {
	_, err := gService.Users.Messages.Trash("me", id).Do()
	return err
}

func UntrashGmailMessage(gService *gmail.Service, id string) error {
	_, err := gService.Users.Messages.Untrash("me", id).Do()
	return err
}

// caution: giving it zero workers would discard the computing at all
func DistributedScan(gService *gmail.Service, gProfile *gmail.Profile, userId int, tasks []*gmail.Message, workers int) []models.Mail {
	channel := make(chan models.Mail, len(tasks)) // buffered ch
	totalTasks := len(tasks)

	if totalTasks == 0 || workers == 0 {
		return nil
	}

	batchSize := totalTasks / min(totalTasks, workers)

	var (
		batchStart = 0
		batchEnd   = batchSize
	)

	for i := range workers {
		// specific check to prevent slicing out of bounds
		if batchStart >= totalTasks {
			break
		}

		// If the calculated end exceeds the total, cap it.
		if batchEnd > totalTasks {
			batchEnd = totalTasks
		}

		go FetchAndParseMail(gService, gProfile, userId, tasks[batchStart:batchEnd], channel, i)

		batchStart = batchEnd
		batchEnd = batchStart + batchSize
	}

	// Handle any remaining items (remainders)
	if batchStart < totalTasks {
		go FetchAndParseMail(gService, gProfile, userId, tasks[batchStart:], channel, workers)
	}

	var mails []models.Mail

	// drain the channel
	// blocking in nature
	for range totalTasks {
		mail := <-channel

		if mail.Id != "" {
			mails = append(mails, mail)
		}
	}

	return mails

}

func FetchAndParseMail(gService *gmail.Service, gProfile *gmail.Profile, userId int, messages []*gmail.Message, ch chan models.Mail, batchNo int) {
	for i, message := range messages {
		message, err := gService.Users.Messages.Get("me", message.Id).Format("metadata").Do()

		// todo: improve this logic
		if err != nil {
			fmt.Println("Failed to fetch message:", message.Id)
			ch <- models.Mail{}
			continue
		}

		senderEmail, err := ExtractSenderAdressFromGmailMessage(message)
		if err != nil {
			fmt.Println("Failed to extract sender email from message", message.Id)
			ch <- models.Mail{}
			continue
		}

		date, err := ExtractDateFromGmailMessage(message)
		if err != nil {
			fmt.Println("Failed to extract date in UTC from message", message.Id)
			ch <- models.Mail{}
			continue
		}

		mail := models.Mail{
			Id:           message.Id,
			AccountEmail: gProfile.EmailAddress,
			SenderEmail:  senderEmail,
			Snippet:      message.Snippet,
			Date:         date,
			UserId:       userId,
		}

		ch <- mail // add this mail to my buffered channel

		// display status
		fmt.Printf("Batch(%d): %d%% - scanned mails - %v\n", batchNo, (i*100)/len(messages), gProfile.EmailAddress)
	}
}
