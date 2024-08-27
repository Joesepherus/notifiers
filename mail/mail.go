package mail

import (
	"net/smtp"
	"os"
)

// sendEmail sends an email notification using Hostinger's SMTP server.
func SendEmail(to, subject, body string) error {
	from := os.Getenv("email") // Replace with your Hostinger email address
	fromName := os.Getenv("name")
	password := os.Getenv("password") // Replace with your Hostinger email password

	// Hostinger SMTP server configuration.
	smtpHost := "smtp.hostinger.com" // Hostinger's SMTP server address
	smtpPort := "587"                // Hostinger's SMTP port (usually 587 for TLS)

	// Message format
	message := []byte(
		"From: " + fromName + " <" + from + ">\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"\r\n" +
			body,
	)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending email.
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		return err
	}
	return nil
}
