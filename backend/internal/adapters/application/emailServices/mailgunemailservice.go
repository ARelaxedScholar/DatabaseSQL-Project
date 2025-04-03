package mailgunemailservice

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

// MailgunEmailService implements the ports.EmailService interface using Mailgun.
type MailgunEmailService struct {
	mg     *mailgun.MailgunImpl
	domain string
	from   string
}

// NewMailgunEmailService creates a new instance of MailgunEmailService.
// domain: your Mailgun domain (e.g., "yourdomain.mailgun.org")
// apiKey: your Mailgun API key
// from: the sender's email address
func NewMailgunEmailService(domain, apiKey, from string) *MailgunEmailService {
	mg := mailgun.NewMailgun(domain, apiKey)
	return &MailgunEmailService{
		mg:     mg,
		domain: domain,
		from:   from,
	}
}

// SendLoginLink sends an email containing a magic login link.
func (s *MailgunEmailService) SendLoginLink(recipient, loginLink string) error {
	subject := "Your Sunflower Booking Magic Login Link"
	text := fmt.Sprintf("Hello,\n\nClick the following link to log in:\n%s\n\nThis link is valid for 10 minutes.", loginLink)

	// Create a new message using Mailgun's API.
	message := s.mg.NewMessage(s.from, subject, text, recipient) // this is technically deprecated, but we are in V4 so it's fine

	// Create a context with a timeout for the API call.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message.
	_, id, err := s.mg.Send(ctx, message)
	if err != nil {
		return err
	}

	// log the Mailgun message ID.
	log.Printf("Mailgun message sent with ID: %s\n", id)
	return nil
}
