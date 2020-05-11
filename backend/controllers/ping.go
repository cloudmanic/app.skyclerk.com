//
// Date: 9/14/2018
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"app.skyclerk.com/backend/services"
)

//
// PingFromClient Collect a ping from the server to know we are alive.
//
func (t *Controller) PingFromClient(c *gin.Context) {
	// // Make sure the UserId is correct.
	// userID := c.MustGet("userId").(uint)
	//
	// // Get the full user
	// _, err := t.db.GetUserById(userID)
	//
	// if err != nil {
	// 	c.JSON(200, gin.H{"status": "logout"})
	// 	return
	// }

	// Get account id
	accountID := uint(c.MustGet("accountId").(int))

	// Get account.
	_, err := t.db.GetAccountById(accountID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found."})
		return
	}

	// Get Billing profile by account id.
	billing, err := t.db.GetBillingByAccountId(accountID)

	if err != nil {
		services.Critical(fmt.Errorf("PingFromClient: Billing account not found. AccountId: %d", accountID))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found (001)."})
		return
	}

	// See if we are Delinquent
	if billing.Status == "Delinquent" {
		c.JSON(200, gin.H{"status": "delinquent"})
		return
	}

	// See if we are Expired
	if billing.Status == "Expired" {
		c.JSON(200, gin.H{"status": "expired"})
		return
	}

	// Return happy JSON
	c.JSON(200, gin.H{"status": "ok"})
}

/* End File */
