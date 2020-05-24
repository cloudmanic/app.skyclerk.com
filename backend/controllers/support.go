//
// Date: 5/24/2019
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"

	"app.skyclerk.com/backend/library/email"
	"app.skyclerk.com/backend/services"
)

const supportEmail = "help@skyclerk.com"

//
// ContactUs will email support after a user requests
//
func (t *Controller) ContactUs(c *gin.Context) {
	// we are only allowed to update certain things.
	body, _ := ioutil.ReadAll(c.Request.Body)
	message := gjson.Get(string(body), "message").String()
	fullName := gjson.Get(string(body), "fullName").String()
	phone := gjson.Get(string(body), "phone").String()
	emailAddress := gjson.Get(string(body), "email").String()

	// Webhook received from fovea
	services.InfoMsg(string(body))

	// Build email.
	subject := fmt.Sprintf("[Website Contact]: New request from %s", fullName)
	html := fmt.Sprintf("<p><b>Name: </b>%s</p> <p><b>Email: </b>%s</p> <p><b>Phone: </b>%s</p> <p><b>Message: </b>%s</p>", fullName, emailAddress, phone, message)

	// Send the email to Support.
	email.Send(supportEmail, emailAddress, subject, html, []string{})

	// Send Slack hook.
	//go slack.Notify("#events", fmt.Sprintf("Skyclerk DoFoveaWebhook Webhook: Email: %s", accountString))

	// Return happy.
	c.JSON(http.StatusNoContent, nil)
}

/* End File */
