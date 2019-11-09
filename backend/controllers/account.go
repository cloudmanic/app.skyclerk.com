//
// Date: 2019-09-16
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"app.skyclerk.com/backend/library/response"
	"app.skyclerk.com/backend/models"
	"app.skyclerk.com/backend/services"
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

//
// DeleteAccount - Delete account.
//
func (t *Controller) DeleteAccount(c *gin.Context) {
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

	// Delete the account.
	t.db.DeleteAccount(account.Id)

	// Get the accounts left for this user.
	u, _ := t.db.GetUserById(uint(userId))

	// Hack for empty accounts.
	if u.Accounts == nil {
		u.Accounts = []models.Account{}
	}

	// Return happy JSON
	c.JSON(200, u.Accounts)
}

//
// NewAccount will create a new account from an account that is already set.
// Sometimes users want to have more accounts under the same billing profile.
//
func (t *Controller) NewAccount(c *gin.Context) {
	// Get body of JSON
	body, _ := ioutil.ReadAll(c.Request.Body)
	name := gjson.Get(string(body), "name").String()

	defer c.Request.Body.Close()

	// Make sure the UserId is correct.
	userID := c.MustGet("userId").(int)

	// Get account id
	accountID := uint(c.MustGet("accountId").(int))

	// Create new account.
	acct := models.Account{
		OwnerId:      uint(userID),
		Name:         strings.Trim(name, " "),
		LastActivity: time.Now(),
	}
	t.db.New().Save(&acct)

	// Add the account look up.
	au := models.AcctToUsers{
		AccountId: acct.Id,
		UserId:    uint(userID),
	}
	t.db.New().Save(&au)

	// Get Billing profile by account id.
	billing, err := t.db.GetBillingByAccountId(accountID)

	if err != nil {
		services.Critical(errors.New(fmt.Sprintf("Billing account not found. AccountId: %d", accountID)))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found (001)."})
		return
	}

	// Update account with billing profile.
	acct.BillingId = billing.Id
	t.db.New().Save(&acct)

	// Get the new account.
	a, err := t.db.GetAccountById(acct.Id)

	if err != nil {
		services.Critical(errors.New(fmt.Sprintf("New account not found. AccountId: %d", accountID)))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found (002)."})
		return
	}

	// Return happy JSON
	c.JSON(200, a)
}

/* End File */
