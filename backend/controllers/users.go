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
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"

	"app.skyclerk.com/backend/emails"
	"app.skyclerk.com/backend/library/email"
	"app.skyclerk.com/backend/library/helpers"
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
// DeleteUser - This does not delete a user just removes from this account.
//
func (t *Controller) DeleteUser(c *gin.Context) {
	// Get the invite id
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	// AccountId.
	accountId := uint(c.MustGet("accountId").(int))

	// Delete AcctToUsers.
	i := models.AcctToUsers{}
	t.db.New().Where("user_id = ? AND acct_id = ?", id, accountId).Delete(&i)

	// Return happy.
	response.RespondDeleted(c, nil)
}

//
// GetInvitedUsers - Return a list of invited users.
//
func (t *Controller) GetInvitedUsers(c *gin.Context) {
	// Get account id
	accountId := uint(c.MustGet("accountId").(int))

	// list of invited users
	list := []models.Invite{}
	t.db.New().Where("expires_at > ? AND account_id = ?", time.Now(), accountId).Find(&list)

	// Return happy.
	response.Results(c, list, nil)
}

//
// DeleteInvite an invite that has been sent out already.
//
func (t *Controller) DeleteInvite(c *gin.Context) {
	// Get the invite id
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	// AccountId.
	accountId := uint(c.MustGet("accountId").(int))

	// Delete invite.
	i := models.Invite{}
	t.db.New().Where("id = ? AND account_id = ?", id, accountId).Delete(&i)

	// Return happy.
	response.RespondDeleted(c, nil)
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
	url := os.Getenv("REGISTER_URL")
	name := me.FirstName + " " + me.LastName
	subject := fmt.Sprintf("%s %s invited you to Skyclerk (%s)", me.FirstName, me.LastName, account.Name)

	// First we check if this user is already in the system
	user, err := t.db.GetUserByEmail(emailAddress)

	// Clean up the message since it is user input
	msg := strings.Replace(template.HTMLEscapeString(message), "\n", "<br>", -1)

	// Not a new user. We should create a create user invite entry.
	if err != nil {
		// Store the user in our invite table
		now := time.Now()
		tExpire := now.Add(time.Hour * 24 * time.Duration(7))

		// Set token
		token := helpers.RandStr(36)

		// Create an invite token
		invite := models.Invite{
			AccountId: account.Id,
			Email:     emailAddress,
			FirstName: firstName,
			LastName:  lastName,
			Message:   msg,
			Token:     token,
			ExpiresAt: tExpire,
		}
		t.db.New().Save(&invite)

		// Update the URL to include our token and names
		url = fmt.Sprintf("%s?email=%s&first=%s&last=%s&token=%s", url, emailAddress, firstName, lastName, token)

		// Send welcome email to user already in the system.
		if flag.Lookup("test.v") != nil {
			email.Send(emailAddress, subject, emails.GetInviteNewUserHTML(name, account.Name, url, invite))
		} else {
			go email.Send(emailAddress, subject, emails.GetInviteNewUserHTML(name, account.Name, url, invite))
		}

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

		// Set URL
		url = os.Getenv("SITE_URL")

		// Create an invite token
		invite := models.Invite{
			AccountId: account.Id,
			Email:     emailAddress,
			FirstName: firstName,
			LastName:  lastName,
			Message:   msg,
		}

		// Send welcome email to user already in the system.
		if flag.Lookup("test.v") != nil {
			email.Send(emailAddress, subject, emails.GetInviteCurrentUserHTML(name, account.Name, url, invite))
		} else {
			go email.Send(emailAddress, subject, emails.GetInviteCurrentUserHTML(name, account.Name, url, invite))
		}

		// Log
		services.InfoMsg(fmt.Sprintf("Current user invited to Skyclerk: AccountId - %d, Email: %s", accountId, emailAddress))
	}

	// Return happy.
	c.JSON(http.StatusNoContent, nil)
}

/* End File */
