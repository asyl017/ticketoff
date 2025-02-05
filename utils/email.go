package utils

import (
	"gopkg.in/gomail.v2"
)

func SendVerificationEmail(to, subject, body string) error {
	from := "sanek.tursumetov@gmail.com"
	password := "sgcutbspyldbjycc"

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, from, password)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
