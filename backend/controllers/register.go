//
// Date: 2019-09-13
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"app.skyclerk.com/backend/library/helpers"
	"app.skyclerk.com/backend/library/realip"
	"app.skyclerk.com/backend/models"
	"app.skyclerk.com/backend/services"
)

const registerStandardErrorMsg = "Something went wrong while logging into your account. Please try again or contact help@skyclerk.com. Sorry for the trouble."

// DoRegister a new account.
func (t *Controller) DoRegister(c *gin.Context) {
	// Set response
	if os.Getenv("APP_ENV") == "local" {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	}

	// Decode json passed in
	decoder := json.NewDecoder(c.Request.Body)

	type RegisterPost struct {
		First    string `json:"first"`
		Last     string `json:"last"`
		Email    string `json:"email"`
		Company  string `json:"company"`
		Password string `json:"password"`
		ClientId string `json:"client_id"`
		Token    string `json:"token"`
	}

	var post RegisterPost

	err := decoder.Decode(&post)

	if err != nil {
		services.Info(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": registerStandardErrorMsg})
		return
	}

	defer c.Request.Body.Close()

	// Validate client id
	app, err := t.db.ValidateClientIdGrantType(post.ClientId, "password")

	if err != nil {
		services.InfoMsg("Register: Invalid client_id")
		c.JSON(http.StatusBadRequest, gin.H{"error": registerStandardErrorMsg})
		return
	}

	// Validate user.
	if err := t.db.ValidateCreateUser(post.First, post.Last, post.Email, false); err != nil {
		// Respond with error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Make sure the password is at least 6 chars long
	if err := t.db.ValidatePassword(post.Password); err != nil {
		// Respond with error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate invite token
	invite := models.Invite{}

	if len(post.Token) > 0 {
		t.db.New().Where("token = ? AND expires_at > ?", post.Token, time.Now()).First(&invite)

		if invite.Id == 0 {
			// Respond with error
			c.JSON(http.StatusBadRequest, gin.H{"error": "Your invite token is not found."})
			return
		}

		// Get the account (this is more or less making sure the account is still there)
		_, err := t.db.GetAccountById(invite.AccountId)

		if err != nil {
			// Respond with error
			c.JSON(http.StatusBadRequest, gin.H{"error": "Your invite token is not found. Unknown account."})
			return
		}
	}

	// Install new user.
	user, err := t.db.CreateUser(post.First, post.Last, post.Email, post.Password, app.Id, c.Request.UserAgent(), realip.RealIP(c.Request))

	if err != nil {
		services.Info(err)

		// Respond with error
		c.JSON(http.StatusBadRequest, gin.H{"error": registerStandardErrorMsg})
		return
	}

	// Figure out name.
	name := post.First + "'s Skyclerk"
	if len(post.Company) > 0 {
		name = post.Company
	}

	// Setup the account
	var acct models.Account

	// Add the account entry
	if len(post.Token) == 0 {
		acct = models.Account{
			OwnerId:      user.Id,
			Name:         name,
			LastActivity: time.Now(),
		}
		t.db.New().Save(&acct)

		// Days to expire
		daysToExpire := helpers.StringToInt(os.Getenv("TRIAL_DAY_COUNT"))

		// Trail expire
		now := time.Now()
		tExpire := now.Add(time.Hour * 24 * time.Duration(daysToExpire))

		// Setup the billing profile for this account.
		bp := models.Billing{
			Status:      "Trial",
			TrialExpire: tExpire,
		}
		t.db.New().Save(&bp)

		// Add billing profile to account.
		acct.BillingId = bp.Id
		t.db.New().Save(&acct)

		// Load default categories.
		t.db.LoadDefaultCategories(acct.Id)
	} else {
		// We know there is no error because this is checked above.
		acct, _ = t.db.GetAccountById(invite.AccountId)

		// Delete invite.
		t.db.New().Delete(&invite)
	}

	// Add the account look up.
	au := models.AcctToUsers{
		AccountId: acct.Id,
		UserId:    user.Id,
	}
	t.db.New().Save(&au)

	// JSON Response
	type Response struct {
		UserId      uint   `json:"user_id"`
		AccessToken string `json:"access_token"`
		AccountId   uint   `json:"account_id"`
	}

	resObj := &Response{
		UserId:      user.Id,
		AccessToken: user.Session.AccessToken,
		AccountId:   acct.Id,
	}

	// Return success json.
	c.JSON(200, resObj)
}

/* End File */
