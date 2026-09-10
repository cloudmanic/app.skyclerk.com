//
// Date: 5/24/2019
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"fmt"
	"html"
	"net/http"

	"github.com/gin-gonic/gin"

	"app.skyclerk.com/backend/library/email"
	"app.skyclerk.com/backend/library/slack"
	"app.skyclerk.com/backend/services"
)

const supportEmail = "help@skyclerk.com"

// Support delivery functions allow tests to assert that rejected requests have no side effects.
var supportSendEmail = email.Send
var supportNotify = slack.Notify

// ContactUs validates a support challenge before sending email or Slack notifications.
// Missing, expired, replayed, or unrelated tokens cannot reach either delivery channel.
func (t *Controller) ContactUs(c *gin.Context) {
	var post struct {
		FullName       string `json:"fullName" binding:"required"`
		Email          string `json:"email" binding:"required,email"`
		Phone          string `json:"phone"`
		Message        string `json:"message" binding:"required"`
		TurnstileToken string `json:"turnstile_token"`
	}
	// Bound this public endpoint's payload before parsing untrusted form fields.
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 65536)
	// Reject malformed JSON and invalid required fields before consuming a challenge.
	if err := c.ShouldBindJSON(&post); err != nil {
		// Return a useful validation error without logging the message or challenge token.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide your name, a valid email address, and a message."})
		return
	}

	// Accept only support challenges issued for the public marketing site's hostnames.
	if !verifyTurnstile(c.Request, post.TurnstileToken, "support_contact", "skyclerk.com", "www.skyclerk.com") {
		// Give the frontend a stable code so it can request a fresh security check.
		c.JSON(http.StatusBadRequest, gin.H{"code": "turnstile_failed", "error": "Please complete the security check and try again."})
		return
	}

	// Build email.
	subject := fmt.Sprintf("[Website Contact]: New request from %s", post.FullName)
	// Keep visitor-provided HTML as text inside the support notification.
	body := fmt.Sprintf("<p><b>Name: </b>%s</p> <p><b>Email: </b>%s</p> <p><b>Phone: </b>%s</p> <p><b>Message: </b>%s</p>", html.EscapeString(post.FullName), html.EscapeString(post.Email), html.EscapeString(post.Phone), html.EscapeString(post.Message))

	// Send the email to Support.
	email.SetNoBccEmail()
	// Send only after verification has succeeded; retain delivery failures for the response.
	err := supportSendEmail(supportEmail, post.Email, subject, body, []string{})
	// Restore the existing BCC behavior for subsequent application emails.
	email.SetBccEmail()
	if err != nil {
		// Log a generic failure without copying the private support message or token.
		services.InfoMsg("Support contact email delivery failed")
		// Preserve the visitor's form so they can retry with a fresh challenge.
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Your request could not be sent. Please try again or email help@skyclerk.com."})
		return
	}

	// Build the notification only after successful email delivery.
	notification := fmt.Sprintf("Skyclerk Website Support Request: Email: %s", post.Email)
	// Send the existing asynchronous Slack notification for verified support requests.
	go supportNotify("#events", notification)

	// Return happy.
	c.JSON(http.StatusNoContent, nil)
}

/* End File */
