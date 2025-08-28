package utils

import (
	"fmt"
	"io"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// EmailMngmnt backed by a real IMAP mailbox
type EmailMngmnt struct {
	Address  string
	Password string
	Server   string // e.g. "imap.gmail.com:993"
}

type Email struct {
	From    string
	Subject string
	Body    string
}

func NewEmailMngmnt(address, password, server string) *EmailMngmnt {
	return &EmailMngmnt{
		Address:  address,
		Password: password,
		Server:   server,
	}
}

func (e *EmailMngmnt) ReadInbox() ([]Email, error) {
	// Connect to server
	c, err := client.DialTLS(e.Server, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	defer c.Logout()

	// Login
	if err := c.Login(e.Address, e.Password); err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("failed to select inbox: %w", err)
	}

	if mbox.Messages == 0 {
		return []Email{}, nil
	}

	// Fetch the last 5 messages
	from := uint32(1)
	if mbox.Messages > 5 {
		from = mbox.Messages - 4
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, mbox.Messages)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		section := &imap.BodySectionName{}
		items := []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}
		done <- c.Fetch(seqset, items, messages)
	}()

	var emails []Email
	for msg := range messages {
		from := msg.Envelope.From[0].MailboxName
		from += "@" + msg.Envelope.From[0].HostName

		// Get the body section
		section := &imap.BodySectionName{}
		r := msg.GetBody(section)
		if r == nil {
			fmt.Println("Server didn't return message body")
			continue
		}

		// Read plaintext body
		b, err := io.ReadAll(r)
		if err != nil {
			fmt.Println("Error reading body:", err)
			continue
		}

		emails = append(emails, Email{
			From:    from,
			Subject: msg.Envelope.Subject,
			Body:    string(b),
		})
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return emails, nil
}
