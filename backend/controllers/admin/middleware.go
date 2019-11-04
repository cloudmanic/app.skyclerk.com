//
// Date: 11/4/2019
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package admin

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"app.skyclerk.com/backend/library/realip"
	"app.skyclerk.com/backend/services"
)

// IPs allowed to access these admin routes
var allowedIps = map[string]bool{
	"127.0.0.1":      true,
	"73.157.194.182": true, // Spicer home
	"35.199.150.98":  true, // vpn.cloudmanic.com
	"208.100.153.45": true, // Bend Condo
}

//
// AuthMiddleware - Here we make sure we passed in a proper Bearer Access Token.
//
func (t *Controller) AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// We only allow this request from a few IP addresses
		if _, ok := allowedIps[realip.RealIP(c.Request)]; !ok {
			services.InfoMsg("UnAuthorization IP address. - " + realip.RealIP(c.Request))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed (#005)"})
			c.AbortWithStatus(401)
			return
		}

		// Set access token and start the auth process
		var access_token = ""

		// Make sure we have a Bearer token.
		auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

		if len(auth) != 2 || auth[0] != "Bearer" {

			// We allow access token from the url
			if os.Getenv("APP_ENV") == "local" {

				access_token = c.Query("access_token")

				if len(access_token) <= 0 {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed (#101)"})
					c.AbortWithStatus(401)
					return
				}

			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed (#001)"})
				c.AbortWithStatus(401)
				return
			}

		} else {
			access_token = auth[1]
		}

		// See if this session is in our db.
		session, err := t.db.GetByAccessToken(access_token)

		if err != nil {
			services.InfoMsg("Access Token Not Found - Unable to Authenticate via HTTP (#002)")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed (#002)"})
			c.AbortWithStatus(401)
			return
		}

		// Get this user is in our db.
		user, err := t.db.GetUserById(session.UserId)

		if err != nil {
			services.InfoMsg("User Not Found - Unable to Authenticate - UserId (HTTP) : " + fmt.Sprint(session.UserId) + " - Session Id : " + fmt.Sprint(session.Id))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed (#003)"})
			c.AbortWithStatus(401)
			return
		}

		// Make sure this is an admin user (Options Cafe Employee)
		if user.Admin != "Yes" {
			services.InfoMsg("User Not Found - Unable to Authenticate - UserId (HTTP) : " + fmt.Sprint(session.UserId) + " - Session Id : " + fmt.Sprint(session.Id))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization Failed (#004)"})
			c.AbortWithStatus(401)
			return
		}

		// Add this user to the context
		c.Set("userId", int(user.Id))

		// CORS for local development.
		if os.Getenv("APP_ENV") == "local" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		// On to next request in the Middleware chain.
		c.Next()
	}
}

/* End File */
