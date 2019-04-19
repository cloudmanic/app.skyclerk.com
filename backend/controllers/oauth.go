//
// Date: 9/8/2017
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cloudmanic/app.skyclerk.com/backend/library/realip"
	"github.com/cloudmanic/app.skyclerk.com/backend/services"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
)

//
// DoOauthToken - Login to account.
//
func (t *Controller) DoOauthToken(c *gin.Context) {
	var username string
	var password string
	var grantType string
	var clientId string

	// A special case to handle clients that do not post in via JSON (looking at you PAW)
	if strings.Contains(c.Request.Header.Get("Content-Type"), "application/x-www-form-urlencoded") {
		username = c.PostForm("username")
		password = c.PostForm("password")
		grantType = c.PostForm("grant_type")
		clientId = c.PostForm("client_id")
	} else {
		body, _ := ioutil.ReadAll(c.Request.Body)
		username = gjson.Get(string(body), "username").String()
		password = gjson.Get(string(body), "password").String()
		grantType = gjson.Get(string(body), "grant_type").String()
		clientId = gjson.Get(string(body), "client_id").String()
	}

	defer c.Request.Body.Close()

	// First we validate the grant type and client id. Make sure this is a known application.
	app, err := t.db.ValidateClientIdGrantType(clientId, grantType)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid client_id or grant type."})
		return
	}

	// Validate user.
	if err := t.db.ValidateUserLogin(username, password); err != nil {

		// Respond with error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Login user in by email and password
	user, sess, err := t.db.LoginUserByEmailPass(username, password, app.Id, c.Request.UserAgent(), realip.RealIP(c.Request))

	if err != nil {
		services.Error(err)

		// Respond with error
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sorry, we could not find your account."})
		return
	}

	//Return success json.
	c.JSON(200, gin.H{"access_token": sess.AccessToken, "user_id": user.Id, "token_type": "bearer"})
}

//
// Logout of account.
//
func (t *Controller) DoLogOut(c *gin.Context) {
	// Search for symbol
	if c.Query("access_token") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Sorry, access_token is required."})
		return
	}

	// Log user out by removing the session
	sess, err := t.db.GetByAccessToken(c.Query("access_token"))

	if err != nil {
		services.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Sorry, we could not find your session."})
		return
	}

	// Delete the session
	t.db.New().Delete(&sess)

	// TODO(spicer): Support passing in ?redirect_url= and then redriect instead of responding with JSON

	// Return success json.
	c.JSON(200, gin.H{"status": "ok"})
}

/* End File */
