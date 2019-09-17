//
// Date: 2019-09-16
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//
// GetAccount returns the account for the logged in user.
//
func (t *Controller) GetAccount(c *gin.Context) {
	// Get account id
	accountId := uint(c.MustGet("accountId").(int))

	// Get account.
	account, err := t.db.GetAccountById(accountId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found."})
		return
	}

	// Return happy JSON
	c.JSON(200, account)
}
