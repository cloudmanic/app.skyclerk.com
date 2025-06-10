//
// Date: 2018-03-21
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-28
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"app.skyclerk.com/backend/library/realip"
	"app.skyclerk.com/backend/models"
	"app.skyclerk.com/backend/services"
)

//
// AuthMiddleware - Here we make sure we passed in a proper Bearer Access Token.
//
func (t *Controller) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate and Set the user.
		user := t.AuthUser(c)

		if user.Id <= 0 {
			return
		}

		// Get the account
		accountID, err := strconv.ParseInt(c.Param("account"), 10, 32)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"system": "Account Not Found - Unable to Authenticate (#005)"}})
			c.AbortWithStatus(401)
			return
		}

		// Make sure this user has this account
		found := false
		for _, row := range user.Accounts {
			if row.Id == uint(accountID) {
				found = true
				break
			}
		}

		if !found {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"system": "Account Not Found - Unable to Authenticate (#006)"}})
			c.AbortWithStatus(401)
			return
		}

		// Set Account
		c.Set("accountId", int(accountID))

		// Set last activity with account.
		go func() {
			account := models.Account{}
			t.db.New().Find(&account, accountID)
			account.LastActivity = time.Now()
			t.db.New().Save(&account)
		}()

		// On to next request in the Middleware chain.
		c.Next()
	}
}

//
// AuthNoAccountMiddleware is used for routes that do not have accounts.
//
func (t *Controller) AuthNoAccountMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate and Set the user.
		t.AuthUser(c)

		// On to next request in the Middleware chain.
		c.Next()
	}
}

//
// AuthUser the user part of the requst.
//
func (t *Controller) AuthUser(c *gin.Context) models.User {
	// Set access token and start the auth process
	var access_token = ""

	// Make sure we have a Bearer token.
	auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Bearer" {

		// We allow access token from the command line (but not during tests)
		if os.Getenv("APP_ENV") == "local" && flag.Lookup("test.v") == nil {

			access_token = c.Query("access_token")

			if len(access_token) <= 0 {
				c.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"system": "Authorization Failed (#001)"}})
				c.AbortWithStatus(401)
				return models.User{}
			}

		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"system": "Authorization Failed (#002)"}})
			c.AbortWithStatus(401)
			return models.User{}
		}

	} else {
		access_token = auth[1]
	}

	// See if this session is in our db.
	session, err := t.db.GetByAccessToken(access_token)

	if err != nil {
		services.LogInfo("Access Token Not Found - Unable to Authenticate via HTTP (#003)")
		c.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"system": "Authorization Failed (#003)"}})
		c.AbortWithStatus(401)
		return models.User{}
	}

	// Get this user is in our db.
	user, err := t.db.GetUserById(session.UserId)

	if err != nil {
		services.LogInfo("User Not Found - Unable to Authenticate - UserId (HTTP) : " + fmt.Sprint(session.UserId) + " - Session Id : " + fmt.Sprint(session.Id))
		c.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"system": "Authorization Failed (#004)"}})
		c.AbortWithStatus(401)
		return models.User{}
	}

	// Log this request into the last_activity col.
	session.LastActivity = time.Now()
	session.LastIpAddress = realip.RealIP(c.Request)
	t.db.New().Save(&session)

	// Add this user to the context
	c.Set("userId", int(user.Id))

	// Return happy
	return user
}

//
// Capture parms out of URL and save them to context. This is useful because we validate integers
// in urls instead of passing what could be a string into an SQL function.
//
func (t *Controller) ParamValidateMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Capture id...
		if len(c.Param("id")) > 0 {

			// Validate input and cast it to a uint.
			id, err := strconv.ParseInt(c.Param("id"), 10, 32)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": gin.H{"system": "The id passed in via the URL is not an integer."}})
				return
			}

			c.Set("id", id)
		}

		// On to the next middleware or the controller.
		c.Next()
	}
}

/* End File */
