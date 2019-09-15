//
// Date: 2019-07-02
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package email

import (
	"errors"
	"os"
	"strconv"

	"app.skyclerk.com/backend/services"
	"gopkg.in/gomail.v2"
	"gopkg.in/mailgun/mailgun-go.v1"
)

var (
	fromEmail = "help@skyclerk.com"
	bccEmail  = "bcc@skyclerk.com"
)

//
// Send - Pass in everything we need to send an email and we send it.
// If we have a SMTP in our configs we use that if not we use
// Mailgun's library for sending mail.
//
func Send(to string, subject string, html string, text string) error {

	// Are we sending as SMTP or via Mailgun? Typically we
	// send as SMTP for local development so we can use Mailhog
	if os.Getenv("MAIL_DRIVER") == "smtp" {
		return SmtpSend(to, subject, html, text)
	}

	// Send via mailgun
	if os.Getenv("MAIL_DRIVER") == "mailgun" {
		return MailgunSend(to, subject, html, text)
	}

	// We should never get here if we are configured correctly.
	var err = errors.New("No mail driver found.")
	services.Info(errors.New(err.Error() + "library/email/Send/Send() - No mail driver found."))
	return err

}

//
// MailgunSend - Send via Mailgun.
//
func MailgunSend(to string, subject string, html string, text string) error {

	// Setup mailgun
	mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_API_KEY"), "")

	// Create message
	message := mailgun.NewMessage("Skyclerk"+"<"+fromEmail+">", subject, text, to)
	message.AddBCC(bccEmail)
	message.SetHtml(html)

	// Send the message
	_, _, err := mg.Send(message)

	if err != nil {
		services.Info(errors.New(err.Error() + "library/email/Send/MailgunSend() - Unable to send email."))
		return err
	}

	// Everything went well!
	return nil

}

//
// SmtpSend - Send as SMTP.
//
func SmtpSend(to string, subject string, html string, text string) error {

	// Setup the email to send.
	m := gomail.NewMessage()
	m.SetHeader("From", fromEmail)
	m.SetHeader("To", to)
	m.SetHeader("Bcc", bccEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", html)
	m.AddAlternative("text/plain", text)

	// Configure the port to be an int.
	port, _ := strconv.ParseInt(os.Getenv("MAIL_PORT"), 10, 64)

	// Make a SMTP connection
	d := gomail.NewDialer(os.Getenv("MAIL_HOST"),
		int(port),
		os.Getenv("MAIL_USERNAME"),
		os.Getenv("MAIL_PASSWORD"))

	// Send Da Email
	if err := d.DialAndSend(m); err != nil {
		services.Info(errors.New(err.Error() + "library/email/Send/SmtpSend() - Unable to send email."))
		return err
	}

	// Everything went well!
	return nil

}

/* End File */
