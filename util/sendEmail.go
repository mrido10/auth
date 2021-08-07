package util

import (
	"auth/config"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

func SendEmail(emailTo []string, emailCc []string, subject string, msg string) error {
	c, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	body := "From: " + c.SendEmail.SENDER_NAME + "\n" +
		"To: " + strings.Join(emailTo, ",") + "\n" +
		"Cc: " + strings.Join(emailCc, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		msg

	auth := smtp.PlainAuth("", c.SendEmail.AUTH_EMAIL, c.SendEmail.AUTH_PASSWORD, c.SendEmail.SMTP_HOST)
	smtpAddr := fmt.Sprintf("%s:%d", c.SendEmail.SMTP_HOST, c.SendEmail.SMTP_PORT)

	err = smtp.SendMail(smtpAddr, auth, c.SendEmail.AUTH_EMAIL, append(emailTo, emailCc...), []byte(body))
	if err != nil {
		return err
	}

	return nil
}
