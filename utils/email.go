package utils

import (
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// EmailMngmnt backed by a real IMAP mailbox
type EmailMngmnt struct {
	Address  string
	Password string
	Server   string // e.g. "imap.gmail.com:993"
}

func NewEmailMngmnt(address, password, server string) *EmailMngmnt {
	return &EmailMngmnt{
		Address:  address,
		Password: password,
		Server:   server,
	}
}

func (e *EmailMngmnt) ReadInbox() ([]string, error) {
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
		return []string{}, nil
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
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	var subjects []string
	for msg := range messages {
		subjects = append(subjects, msg.Envelope.Subject)
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return subjects, nil
}
