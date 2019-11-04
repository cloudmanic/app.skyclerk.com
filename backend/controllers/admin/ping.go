//
// Date: 11/3/2018
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package admin

import (
	"github.com/gin-gonic/gin"
)

//
// PingFromServer - Collect a ping from the server to know we are alive.
// Mainly used as a way to see if a user is logged in an allowed
// to access the admin side of the application.
//
func (t *Controller) PingFromServer(c *gin.Context) {

	// Make sure the UserId is correct.
	//userId := c.MustGet("userId").(uint)

	// Return happy JSON
	c.JSON(200, gin.H{"status": "ok"})
}

/* End File */
