package models

import (
	"fmt"
	"time"
)

type Mail struct {
	Id           string
	UserId       int
	AccountEmail string
	SenderEmail  string
	Snippet      string
	Date         time.Time
}

func DeleteMailsByUserId(userId int) error {
	_, err := DB.Exec("DELETE FROM mails WHERE user_id=?", userId)
	return err
}

func DeleteMail(mailId string) error {
	_, err := DB.Exec("DELETE FROM mails WHERE id=?", mailId)
	return err
}

func AddMails(mails []Mail) error {
	// prepare the transaction for the db
	tx, err := DB.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO mails(id, account_email, sender_email, snippet, date, user_id) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, mail := range mails {
		if _, err := stmt.Exec(mail.Id, mail.AccountEmail, mail.SenderEmail, mail.Snippet, mail.Date.Format(time.DateTime), mail.UserId); err != nil {
			tx.Rollback()
			return err
		}
	}

	// commit the transaction
	return tx.Commit()
}

func ListMailsBySender(senderEmail string, userId int) ([]Mail, error) {
	rows, err := DB.Query("SELECT id, account_email, sender_email, snippet FROM mails WHERE user_id=? AND sender_email=?", userId, senderEmail)

	if err != nil {
		return []Mail{}, err
	}

	mails := []Mail{}

	for rows.Next() {
		data := Mail{}
		err := rows.Scan(&data.Id, &data.AccountEmail, &data.SenderEmail, &data.Snippet)

		if err != nil {
			fmt.Println("skipping read from tables `mails`", err)
			continue
		}
		mails = append(mails, data)
	}

	return mails, nil
}

func ListMails(accountEmail string) ([]Mail, error) {
	rows, err := DB.Query("SELECT id, account_email, sender_email, snippet FROM mails WHERE account_email=?", accountEmail)

	if err != nil {
		return []Mail{}, err
	}

	mails := []Mail{}

	for rows.Next() {
		data := Mail{}
		err := rows.Scan(&data.Id, &data.AccountEmail, &data.SenderEmail, &data.Snippet)

		if err != nil {
			fmt.Println("skipping read from tables `mails`", err)
			continue
		}
		mails = append(mails, data)
	}

	return mails, nil
}

type MailFolder struct {
	SenderEmail string
	Count       int
}

func ListMailFolders(userId int) ([]MailFolder, error) {
	rows, err := DB.Query("SELECT sender_email, COUNT(sender_email) FROM mails WHERE user_id=? GROUP BY sender_email ORDER BY COUNT(sender_email) DESC", userId)

	if err != nil {
		return []MailFolder{}, err
	}

	folders := []MailFolder{}

	for rows.Next() {
		data := MailFolder{}
		err := rows.Scan(&data.SenderEmail, &data.Count)
		if err != nil {
			fmt.Println("skipping read from tables `mails`", err)
			continue
		}

		folders = append(folders, data)
	}

	return folders, nil
}
