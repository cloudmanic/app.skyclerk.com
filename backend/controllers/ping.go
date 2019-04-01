//
// Date: 9/14/2018
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"github.com/gin-gonic/gin"
)

//
// Collect a ping from the server to know we are alive.
//
func (t *Controller) PingFromServer(c *gin.Context) {

	// // Make sure the UserId is correct.
	// userId := c.MustGet("userId").(uint)
	//
	// // Get the full user
	// user, err := t.DB.GetUserById(userId)
	//
	// if err != nil {
	// 	c.JSON(200, gin.H{"status": "logout"})
	// 	return
	// }
	//
	// // See if we are Delinquent
	// if user.Status == "Delinquent" {
	// 	c.JSON(200, gin.H{"status": "delinquent"})
	// 	return
	// }
	//
	// // See if we are Expired
	// if user.Status == "Expired" {
	// 	c.JSON(200, gin.H{"status": "expired"})
	// 	return
	// }

	// Return happy JSON
	c.JSON(200, gin.H{"status": "ok"})
}

/* End File */
