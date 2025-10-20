package jobs

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
)

// SendEmail sends an email with the given subject and body to recipient.
// Uses Gmail SMTP with app password from environment variables.
func SendEmail(to, subject, body string) error {
	sender := os.Getenv("EMAIL_SENDER")
	password := os.Getenv("EMAIL_PASSWORD")
	if sender == "" || password == "" {
		return fmt.Errorf("EMAIL_SENDER or EMAIL_PASSWORD not set")
	}

	host := "smtp.gmail.com"
	port := 587
	auth := smtp.PlainAuth("", sender, password, host)

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n\r\n%s\r\n", to, subject, body))

	// Upgrade to TLS after connecting on 587
	client, err := smtp.Dial(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	defer client.Quit()

	if err := client.StartTLS(&tls.Config{ServerName: host}); err != nil {
		return err
	}
	if err := client.Auth(auth); err != nil {
		return err
	}
	if err := client.Mail(sender); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}
	wc, err := client.Data()
	if err != nil {
		return err
	}
	defer wc.Close()

	if _, err := wc.Write(msg); err != nil {
		return err
	}
	return nil
}
