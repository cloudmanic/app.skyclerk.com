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
	"os"
	"strings"
	"time"

	"app.skyclerk.com/backend/library/response"
	"app.skyclerk.com/backend/library/stripe"
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

	// Put default categories in.
	t.db.LoadDefaultCategories(acct.Id)

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

//
// NewStripeToken will update the credit card on file with stripe.
//
func (t *Controller) NewStripeToken(c *gin.Context) {
	// we are only allowed to update certain things.
	body, _ := ioutil.ReadAll(c.Request.Body)
	plan := gjson.Get(string(body), "plan").String()
	token := gjson.Get(string(body), "token").String()

	// Make sure the UserId is correct.
	userID := c.MustGet("userId").(int)

	// Get account id
	accountID := uint(c.MustGet("accountId").(int))

	// Get account.
	account, err := t.db.GetAccountById(accountID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found."})
		return
	}

	// We must be the account owner to proceed
	if account.OwnerId != uint(userID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You must be the account owner."})
		return
	}

	// Get the user attached to this account.
	user, err := t.db.GetUserById(uint(userID))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not user."})
		return
	}

	// Make sure we have a token
	if len(token) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You must include a token."})
		return
	}

	// Get Billing profile by account id.
	billing, err := t.db.GetBillingByAccountId(accountID)

	if err != nil {
		services.Critical(errors.New(fmt.Sprintf("NewStripeToken: Billing account not found. AccountId: %d", accountID)))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found (001)."})
		return
	}

	// Create a new stripe customer.
	custID := billing.StripeCustomer

	if len(billing.StripeCustomer) == 0 {
		custID, err = stripe.AddCustomer(user.FirstName, user.LastName, user.Email, int(accountID))

		if err != nil {
			services.Critical(errors.New(fmt.Sprintf("Error with Stripe new account. AccountId: %d - %s", accountID, err.Error())))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown error. Please contact help@skyclerk.com."})
			return
		}
	}

	// Delete all cards on file with stripe.
	cards, err := stripe.ListAllCreditCards(custID)

	if err != nil {
		services.Critical(errors.New(fmt.Sprintf("Error with Stripe listng cards. AccountId: %d - %s", accountID, err.Error())))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown error. Please contact help@skyclerk.com."})
		return
	}

	for _, row := range cards {
		err = stripe.DeleteCreditCard(custID, row)

		if err != nil {
			services.Critical(errors.New(fmt.Sprintf("Error with Stripe delete card. AccountId: %d - %s", accountID, err.Error())))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown error. Please contact help@skyclerk.com."})
			return
		}
	}

	// Add card to customer.
	_, err = stripe.AddCreditCardByToken(custID, token)

	if err != nil {
		services.Critical(errors.New(fmt.Sprintf("Error with Stripe new card. AccountId: %d - %s", accountID, err.Error())))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown error. Please contact help@skyclerk.com."})
		return
	}

	// Get planID yearly or monthly
	planID := os.Getenv("STRIPE_MONTHLY_PLAN")
	billing.Subscription = "Monthly"
	if plan == "yearly" {
		planID = os.Getenv("STRIPE_YEARLY_PLAN")
		billing.Subscription = "Yearly"
	}

	// Create a new subscription if need be
	subID := billing.StripeSubscription

	if len(billing.StripeSubscription) == 0 {
		subID, err = stripe.AddSubscription(custID, planID, "", false)

		if err != nil {
			services.Critical(errors.New(fmt.Sprintf("Error with Stripe new subscription. AccountId: %d - %s", accountID, err.Error())))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown error. Please contact help@skyclerk.com."})
			return
		}
	}

	// Update billing profile.
	billing.Status = "Active"
	billing.StripeCustomer = custID
	billing.StripeSubscription = subID
	t.db.New().Save(&billing)

	// Return happy.
	c.JSON(http.StatusNoContent, nil)
}

//
// ChangeSubscription will change the plan the user is on.
//
func (t *Controller) ChangeSubscription(c *gin.Context) {
	// we are only allowed to update certain things.
	body, _ := ioutil.ReadAll(c.Request.Body)
	plan := gjson.Get(string(body), "plan").String()

	// Make sure the UserId is correct.
	userID := c.MustGet("userId").(int)

	// Get account id
	accountID := uint(c.MustGet("accountId").(int))

	// Get account.
	account, err := t.db.GetAccountById(accountID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found."})
		return
	}

	// We must be the account owner to proceed
	if account.OwnerId != uint(userID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You must be the account owner."})
		return
	}

	// Get Billing profile by account id.
	billing, err := t.db.GetBillingByAccountId(accountID)

	if err != nil {
		services.Critical(errors.New(fmt.Sprintf("ChangeSubscription: Billing account not found. AccountId: %d", accountID)))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found (001)."})
		return
	}

	// Make sure we currently have a subscription
	if (len(billing.StripeCustomer) == 0) || (len(billing.StripeSubscription) == 0) {
		services.Critical(errors.New(fmt.Sprintf("ChangeSubscription: Billing account not found. AccountId: %d", accountID)))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found (002)."})
		return
	}

	// Get planID yearly or monthly
	planID := os.Getenv("STRIPE_MONTHLY_PLAN")
	billing.Subscription = "Monthly"
	if plan == "Yearly" {
		planID = os.Getenv("STRIPE_YEARLY_PLAN")
		billing.Subscription = "Yearly"
	}

	// Update subscription
	subID, err := stripe.UpdateSubscription(billing.StripeSubscription, planID)

	if err != nil {
		services.Critical(errors.New(fmt.Sprintf("Error with Stripe update subscription. AccountId: %d - %s", accountID, err.Error())))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown error. Please contact help@skyclerk.com."})
		return
	}

	// Update billing profile.
	billing.Status = "Active"
	billing.StripeSubscription = subID
	t.db.New().Save(&billing)

	// Return happy.
	c.JSON(http.StatusNoContent, nil)
}

/* End File */
