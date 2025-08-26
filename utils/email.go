package utils

import (
	"github.com/z11i/onesecmail"
)

type EmailMngmnt struct {
	username string
}

type Email struct {
	From        string
	Subject     string
	Date        string
	Attachments []string
	Body        *string
}

func NewEmailMngmnt(username string) *EmailMngmnt {
	return &EmailMngmnt{username: username}
}

func (e *EmailMngmnt) ReadInbox() ([]Email, error) {
	mailbox, err := onesecmail.NewMailbox(e.username, "1secmail.com", nil)
	if err != nil {
		return nil, err
	}

	items, err := mailbox.CheckInbox()
	if err != nil {
		return nil, err
	}

	var emails []Email

	for _, item := range items {
		m, err := mailbox.ReadMessage(item.ID)
		if err != nil {
			continue
		}
		emails = append(emails, Email{
			From:    m.From,
			Subject: m.Subject,
			Date:    m.Date,
			Body:    m.Body,
		})
	}

	return emails, nil
}
