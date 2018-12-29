//
// Date: 2018-03-21
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-28
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

//
// Here we make sure we passed in a proper Bearer Access Token.
//
func (t *Controller) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// // Set access token and start the auth process
		// var access_token = ""

		// // Make sure we have a Bearer token.
		// auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		// if len(auth) != 2 || auth[0] != "Bearer" {

		// 	// We allow access token from the command line
		// 	if os.Getenv("APP_ENV") == "local" {

		// 		access_token = c.Query("access_token")

		// 		if len(access_token) <= 0 {
		// 			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed (#101)"})
		// 			c.AbortWithStatus(401)
		// 			return
		// 		}

		// 	} else {
		// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed (#001)"})
		// 		c.AbortWithStatus(401)
		// 		return
		// 	}

		// } else {
		// 	access_token = auth[1]
		// }

		// // See if this session is in our db.
		// session, err := t.db.GetByAccessToken(access_token)

		// if err != nil {
		// 	services.Critical("Access Token Not Found - Unable to Authenticate via HTTP (#002)")
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed (#002)"})
		// 	c.AbortWithStatus(401)
		// 	return
		// }

		// // Get this user is in our db.
		// user, err := t.db.GetUserById(session.UserId)

		// if err != nil {
		// 	services.Critical("User Not Found - Unable to Authenticate - UserId (HTTP) : " + fmt.Sprint(session.UserId) + " - Session Id : " + fmt.Sprint(session.Id))
		// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed (#003)"})
		// 	c.AbortWithStatus(401)
		// 	return
		// }

		// // Log this request into the last_activity col.
		// session.LastActivity = time.Now()
		// session.LastIpAddress = realip.RealIP(c.Request)
		// t.db.UpdateSession(&session)

		// Add this user to the context
		c.Set("userId", uint(109))

		// TODO: verify this user is allowed to access this account.
		// This is set in ParamValidateMiddleware.
		// c.MustGet("account").(int)

		// CORS for local deve opment.
		if os.Getenv("APP_ENV") == "local" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		// On to next request in the Middleware chain.
		c.Next()
	}
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
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "The id passed in via the URL is not an integer."})
				return
			}

			c.Set("id", id)
		}

		// Capture account...
		if len(c.Param("account")) > 0 {

			// Validate input and cast it to a uint.
			account, err := strconv.ParseInt(c.Param("account"), 10, 32)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "The account passed in via the URL is not an integer."})
				return
			}

			c.Set("account", account)
		}

		// On to the next middleware or the controller.
		c.Next()
	}
}

/* End File */
