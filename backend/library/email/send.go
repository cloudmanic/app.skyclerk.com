//
// Date: 2019-07-02
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package email

import (
	"bufio"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/keighl/postmark"
	"gopkg.in/gomail.v2"
	"gopkg.in/mailgun/mailgun-go.v1"

	"app.skyclerk.com/backend/library/files"
	"app.skyclerk.com/backend/library/html2text"
	"app.skyclerk.com/backend/services"
)

var (
	fromName  = "Skyclerk"
	fromEmail = "help@skyclerk.com"
	bccEmail  = "bcc@skyclerk.com"
)

//
// Send - Pass in everything we need to send an email and we send it.
// If we have a SMTP in our configs we use that if not we use
// Mailgun's library for sending mail. Attachments are an array of local file paths.
//
func Send(to string, replyTo string, subject string, html string, attachments []string) error {
	// Setup text email
	text, err := html2text.FromString(html, html2text.Options{})

	if err != nil {
		return err
	}

	// Override reply to if empty
	if replyTo == "" {
		replyTo = fromEmail
	}

	// Are we sending as SMTP or via Mailgun? Typically we
	// send as SMTP for local development so we can use Mailhog
	if os.Getenv("MAIL_DRIVER") == "smtp" {
		return SMTPSend(to, replyTo, subject, html, text, attachments)
	}

	// Send via mailgun
	if os.Getenv("MAIL_DRIVER") == "mailgun" {
		return MailgunSend(to, replyTo, subject, html, text, attachments)
	}

	// Send via postmark
	if os.Getenv("MAIL_DRIVER") == "postmark" {
		return PostmarkSend(to, replyTo, subject, html, text, attachments)
	}

	// We should never get here if we are configured correctly.
	err = errors.New("No mail driver found.")
	services.Info(errors.New(err.Error() + "library/email/Send/Send() - No mail driver found."))
	return err

}

//
// PostmarkSend will send via postmark.
//
func PostmarkSend(to string, replyTo string, subject string, html string, text string, attachments []string) error {
	// Setup postmark
	client := postmark.NewClient(os.Getenv("POSTMARK_SERVER_KEY"), os.Getenv("POSTMARK_ACCOUNT_KEY"))

	email := postmark.Email{
		From:       fromEmail,
		To:         to,
		ReplyTo:    replyTo,
		Bcc:        bccEmail,
		Subject:    subject,
		HtmlBody:   html,
		TextBody:   text,
		TrackOpens: true,
	}

	// Include any attachements.
	for _, row := range attachments {
		// Open file on disk.
		f, _ := os.Open(row)

		// Read entire JPG into byte slice.
		reader := bufio.NewReader(f)
		content, _ := ioutil.ReadAll(reader)

		// Encode as base64.
		encoded := base64.StdEncoding.EncodeToString(content)

		// Get the file type
		contentType, _, err := files.FileContentTypeWithError(row)

		if err != nil {
			services.Info(errors.New(err.Error() + "library/email/Send/PostmarkSend() - Unable to get contentType."))
			return err
		}

		// Build atachement object.
		a := postmark.Attachment{
			Name:        path.Base(row),
			Content:     encoded,
			ContentType: contentType,
		}

		email.Attachments = append(email.Attachments, a)

		f.Close()
	}

	// Send the email.
	_, err := client.SendEmail(email)

	if err != nil {
		services.Info(errors.New(err.Error() + "library/email/Send/PostmarkSend() - Unable to send email."))
		return err
	}

	// Everything went well!
	return nil
}

//
// MailgunSend - Send via Mailgun.
//
func MailgunSend(to string, replyTo string, subject string, html string, text string, attachments []string) error {
	// Setup mailgun
	mg := mailgun.NewMailgun(os.Getenv("MAILGUN_DOMAIN"), os.Getenv("MAILGUN_API_KEY"), "")

	// Create message
	message := mailgun.NewMessage(fromName+"<"+fromEmail+">", subject, text, to)
	message.AddBCC(bccEmail)
	message.SetHtml(html)
	message.SetReplyTo(replyTo)

	// Include any attachements.
	for _, row := range attachments {
		message.AddAttachment(row)
	}

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
// SMTPSend - Send as SMTP.
//
func SMTPSend(to string, replyTo string, subject string, html string, text string, attachments []string) error {
	// Setup the email to send.
	m := gomail.NewMessage()
	m.SetHeader("From", fromEmail)
	m.SetHeader("To", to)
	m.SetHeader("ReplyTo", replyTo)
	m.SetHeader("Bcc", bccEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", html)
	m.AddAlternative("text/plain", text)

	// Include any attachements.
	for _, row := range attachments {
		m.Attach(row)
	}

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
