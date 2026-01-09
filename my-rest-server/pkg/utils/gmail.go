package utils

import (
	"context"
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

	gService.Users.Messages.List("me").Pages(
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
