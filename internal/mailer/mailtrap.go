package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	gomail "gopkg.in/mail.v2"
)

type MailTrapMailer struct {
	fromEmail string
	apiKey    string
}

func NewMailTrapMailer(apiKey, fromEmail string) (*MailTrapMailer, error) {
	if apiKey == "" {
		return nil, errors.New("api key is required")
	}

	return &MailTrapMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
	}, nil
}

func (m *MailTrapMailer) Send(templateFile, username, email string, data any, isSandbox bool) (int, error) {
	// template parsing and building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)

	if err != nil {
		return -1, err
	}

	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return -1, err
	}

	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return -1, err
	}

	message := gomail.NewMessage()

	message.SetHeader("From", m.fromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject.String())

	message.AddAlternative("text/html", body.String())

	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", m.apiKey)

	var retryError error
	for i := 0; i < maxRetries; i++ {
		retryError = dialer.DialAndSend(message)

		if retryError != nil {
			log.Printf("Failed to send email to %s, attempt %d of %d", email, i+1, maxRetries)

			//exponential backoff
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}

		log.Printf("Email sent with status code %v", http.StatusOK)
		return http.StatusOK, nil
	}

	return -1, fmt.Errorf("failed to send email after %d attempts, error: %v", maxRetries, retryError)
}
