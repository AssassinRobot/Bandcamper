package utils

type EmailMngmnt struct {
	address string
}

func NewEmailMngmnt(address string) *EmailMngmnt {
	return &EmailMngmnt{address: address}
}

func (e *EmailMngmnt) ReadInbox() ([]string, error) {
	// Placeholder for reading emails from the inbox
	// In a real implementation, this would connect to an email server
	// and fetch emails. Here we return a mock list of email subjects.
	emails := []string{
		"Welcome to Bandcamper!",
		"Your download is ready",
		"Weekly newsletter",
	}
	return emails, nil
}
