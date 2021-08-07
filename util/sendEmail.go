package util

import (
	"auth/config"
	"fmt"
	"log"

	"gopkg.in/gomail.v2"
)

func SendEmail(emailTo string, emailCc string, subject string, msg string, template string) error {
	c, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", c.SendEmail.SENDER_NAME+"<"+c.SendEmail.AUTH_EMAIL+">")
	mailer.SetHeader("To", emailTo)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", msg)

	dialer := gomail.NewDialer(
		c.SendEmail.SMTP_HOST,
		c.SendEmail.SMTP_PORT,
		c.SendEmail.AUTH_EMAIL,
		c.SendEmail.AUTH_PASSWORD,
	)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}
