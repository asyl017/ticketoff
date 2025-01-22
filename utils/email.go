package utils

import "net/smtp"

func SendEmail(to, subject, body string) error {
	from := "your-email@example.com"
	password := "your-email-password"

	// Set up authentication information.
	auth := smtp.PlainAuth("", from, password, "smtp.example.com")

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body

	err := smtp.SendMail("smtp.example.com:587", auth, from, []string{to}, []byte(msg))
	return err
}
