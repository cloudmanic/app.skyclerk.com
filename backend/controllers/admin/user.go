//
// Date: 11/3/2018
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//
// GetUser - Return a single user.
//
func (t *Controller) GetUser(c *gin.Context) {
	// Set user id
	var userID = c.MustGet("id").(int64)

	// Get the full user
	user, err := t.db.GetUserById(uint(userID))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found."})
		return
	}

	// Return happy JSON
	c.JSON(200, user)
}

/* End File */
