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

	"github.com/cloudmanic/app.options.cafe/backend/library/helpers"
	"github.com/gin-gonic/gin"

	"app.skyclerk.com/backend/library/realip"
	"app.skyclerk.com/backend/models"
	"app.skyclerk.com/backend/services"
)

const registerStandardErrorMsg = "Something went wrong while logging into your account. Please try again or contact help@options.cafe. Sorry for the trouble."

//
// DoRegister a new account.
//
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
		Password string `json:"password"`
		ClientId string `json:"client_id"`
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

	// Install new user.
	user, err := t.db.CreateUser(post.First, post.Last, post.Email, post.Password, app.Id, c.Request.UserAgent(), realip.RealIP(c.Request))

	if err != nil {
		services.Info(err)

		// Respond with error
		c.JSON(http.StatusBadRequest, gin.H{"error": registerStandardErrorMsg})
		return
	}

	// Days to expire
	daysToExpire := helpers.StringToInt(os.Getenv("TRIAL_DAY_COUNT"))

	// Trail expire
	now := time.Now()
	tExpire := now.Add(time.Hour * 24 * time.Duration(daysToExpire))

	// Add the account entry
	acct := models.Account{
		OwnerId:      user.Id,
		Name:         post.First + "'s Skyclerk",
		Status:       "Trial",
		LastActivity: time.Now(),
		SignupIp:     realip.RealIP(c.Request),
		TrialExpire:  tExpire,
	}
	t.db.New().Save(&acct)

	// Add the account look up.
	au := models.AcctToUsers{
		AcctId: acct.Id,
		UserId: user.Id,
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
