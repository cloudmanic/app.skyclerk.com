//
// Date: 2019-07-02
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"

	"app.skyclerk.com/backend/emails"
	"app.skyclerk.com/backend/library/email"
	"app.skyclerk.com/backend/library/response"
	"app.skyclerk.com/backend/models"
	"app.skyclerk.com/backend/services"
)

//
// GetUsers - Return a list of users. We limit to 500 mainly so we do not overload the
// system, but enough so the front-end does not have to page
//
func (t *Controller) GetUsers(c *gin.Context) {
	// Get account id
	accountId := uint(c.MustGet("accountId").(int))

	// Query and get users
	users := t.db.GetUsersByAccount(accountId)

	// Return happy.
	response.Results(c, users, nil)
}

//
// InviteUser - Invites a new user to an account. If the user is new they will create an account.
// If not we will simply add them.
//
func (t *Controller) InviteUser(c *gin.Context) {
	// Data posted in.
	var firstName string
	var lastName string
	var emailAddress string
	var message string

	// Get account id
	accountId := uint(c.MustGet("accountId").(int))

	// Make sure the UserId is correct.
	userId := uint(c.MustGet("userId").(int))

	// Get the full user
	me, err := t.db.GetUserById(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found."})
		return
	}

	// Get account.
	account, err := t.db.GetAccountById(accountId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found."})
		return
	}

	// Parse JSON posted in.
	body, _ := ioutil.ReadAll(c.Request.Body)
	firstName = gjson.Get(string(body), "first_name").String()
	lastName = gjson.Get(string(body), "last_name").String()
	emailAddress = gjson.Get(string(body), "email").String()
	message = gjson.Get(string(body), "message").String()

	// Poor man's validation
	if (len(firstName) == 0) || (len(lastName) == 0) || (len(emailAddress) == 0) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The following fields are required: first_name, last_name, email."})
		return
	}

	// Make sure we have a valid email address
	err = t.db.ValidateEmailAddress(emailAddress)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email address is not valid."})
		return
	}

	// Setup the basic vars for sending emails.
	url := os.Getenv("SITE_URL")
	name := me.FirstName + " " + me.LastName
	subject := fmt.Sprintf("%s %s invited you to Skyclerk (%s)", me.FirstName, me.LastName, account.Name)

	// First we check if this user is already in the system
	user, err := t.db.GetUserByEmail(emailAddress)

	// Not a new user. We should create a create user invite entry.
	if err != nil {
		// Store the user in our invite table
		now := time.Now()
		tExpire := now.Add(time.Hour * 24 * time.Duration(7))

		// Create an invite token
		invite := models.Invite{
			AccountId: account.Id,
			Email:     emailAddress,
			FirstName: firstName,
			LastName:  lastName,
			Message:   message,
			Token:     "abc123",
			ExpiresAt: tExpire,
		}
		t.db.New().Save(&invite)

		// // Send a welcome email to user. This is different than the welcome email below.
		// html := emails.GetInviteCurrentUserHTML(name, account.Name, url)
		// text := emails.GetInviteCurrentUserText(name, account.Name, url)
		//
		// // Send welcome email to user already in the system.
		// if flag.Lookup("test.v") != nil {
		// 	email.Send(emailAddress, subject, html, text)
		// } else {
		// 	go email.Send(emailAddress, subject, html, text)
		// }

		// Log
		services.InfoMsg(fmt.Sprintf("New user invited to Skyclerk: AccountId - %d, Email: %s", accountId, emailAddress))
	} else { // current user
		// Validate this user is not already part of the account.
		u := models.AcctToUsers{}
		t.db.New().Where("acct_id = ? AND  user_id = ?", account.Id, user.Id).First(&u)

		if u.Id > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User is already part of this account."})
			return
		}

		// User already in the system let's just assign them to the account.
		t.db.New().Save(&models.AcctToUsers{
			AcctId: accountId,
			UserId: user.Id,
		})

		// Setup emails to send for current user
		html := emails.GetInviteCurrentUserHTML(name, account.Name, url)
		text := emails.GetInviteCurrentUserText(name, account.Name, url)

		// Send welcome email to user already in the system.
		if flag.Lookup("test.v") != nil {
			email.Send(emailAddress, subject, html, text)
		} else {
			go email.Send(emailAddress, subject, html, text)
		}

		// Log
		services.InfoMsg(fmt.Sprintf("Current user invited to Skyclerk: AccountId - %d, Email: %s", accountId, emailAddress))
	}

	// Return happy.
	c.JSON(http.StatusNoContent, nil)
}

/* End File */
