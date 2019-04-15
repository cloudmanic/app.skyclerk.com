//
// Date: 4/14/2019
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//
// GetMe returns all the data related to a user.
//
func (t *Controller) GetMe(c *gin.Context) {
	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(uint)

	// Get the full user
	user, err := t.db.GetUserById(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found."})
		return
	}

	// Return happy JSON
	c.JSON(200, user)
}
