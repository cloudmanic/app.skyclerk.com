//
// Date: 2019-09-16
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"io/ioutil"
	"net/http"
	"strings"

	"app.skyclerk.com/backend/library/response"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
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

//
// UpdateAccount - Update a account.
//
func (t *Controller) UpdateAccount(c *gin.Context) {
	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(int)

	// Get account id
	accountId := uint(c.MustGet("accountId").(int))

	// Get account.
	account, err := t.db.GetAccountById(accountId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found."})
		return
	}

	// We must be the account owner to proceed
	if account.OwnerId != uint(userId) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You must be the account owner."})
		return
	}

	// we are only allowed to update certain things.
	body, _ := ioutil.ReadAll(c.Request.Body)
	name := gjson.Get(string(body), "name").String()
	currency := gjson.Get(string(body), "currency").String()
	locale := gjson.Get(string(body), "locale").String()
	ownerId := gjson.Get(string(body), "owner_id").Int()

	// Update object
	account.OwnerId = uint(ownerId)
	account.Name = strings.Trim(name, " ")
	account.Locale = strings.Trim(locale, " ")
	account.Currency = strings.Trim(currency, " ")

	// Valdate the data in the model.
	err2 := account.Validate(t.db, "update", uint(userId), accountId, account.Id)

	// If we had validation errors return them and do no more.
	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err2})
		return
	}

	// Update Account
	t.db.New().Save(&account)

	// Get account - Refresh the data.
	accountNew, err := t.db.GetAccountById(accountId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found."})
		return
	}

	// Return happy.
	response.RespondUpdated(c, accountNew, nil)
}

//
// ClearAccount - Clear account.
//
func (t *Controller) ClearAccount(c *gin.Context) {
	// Make sure the UserId is correct.
	userId := c.MustGet("userId").(int)

	// Get account id
	accountId := uint(c.MustGet("accountId").(int))

	// Get account.
	account, err := t.db.GetAccountById(accountId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found."})
		return
	}

	// We must be the account owner to proceed
	if account.OwnerId != uint(userId) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You must be the account owner."})
		return
	}

	// Clear the account.
	t.db.ClearAccount(account.Id)

	// Put default categories back in.
	t.db.LoadDefaultCategories(account.Id)

	// Return happy.
	c.JSON(http.StatusNoContent, nil)
}

/* End File */
