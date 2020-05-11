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
		stripe.DeleteCustomer(custID)
		stripeErr := gjson.Get(err.Error(), "message").String()
		services.Critical(fmt.Errorf("Error with Stripe new card. AccountId: %d - %s", accountID, err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": stripeErr})
		return
	}

	// Get planID yearly or monthly
	planID := os.Getenv("STRIPE_MONTHLY_PLAN")
	billing.Subscription = "Monthly"
	if plan == "Yearly" {
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
		services.Critical(fmt.Errorf("Error with Stripe update subscription. AccountId: %d - %s", accountID, err.Error()))
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

//
// GetBilling will return the billing informaton for the account.
//
func (t *Controller) GetBilling(c *gin.Context) {
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
		services.Critical(fmt.Errorf("GetBilling: Billing account not found. AccountId: %d", accountID))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found (001)."})
		return
	}

	// If we have a stripe user
	if len(billing.StripeCustomer) == 0 {
		// No stripe user.
		c.JSON(200, billing)
		return
	}

	// Get stripe customer
	stripeCustomer, err := stripe.GetCustomer(billing.StripeCustomer)

	if err != nil {
		// No stripe user.
		c.JSON(200, billing)
		return
	}

	// Statement stuff.
	billing.CurrentPeriodStart = time.Unix(stripeCustomer.Subscriptions.Data[0].CurrentPeriodStart, 0)
	billing.CurrentPeriodEnd = time.Unix(stripeCustomer.Subscriptions.Data[0].CurrentPeriodEnd, 0)

	// Do we have a credit card on file
	if stripeCustomer.Sources.ListMeta.TotalCount > 0 {
		billing.CardBrand = string(stripeCustomer.Sources.Data[0].Card.Brand)
		billing.CardLast4 = stripeCustomer.Sources.Data[0].Card.Last4
		billing.CardExpMonth = int(stripeCustomer.Sources.Data[0].Card.ExpMonth)
		billing.CardExpYear = int(stripeCustomer.Sources.Data[0].Card.ExpYear)
	}

	// Return happy JSON
	c.JSON(200, billing)
}

//
// GetBillingHistory will return the billing history
//
func (t *Controller) GetBillingHistory(c *gin.Context) {
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
		services.Critical(fmt.Errorf("GetBillingHistory: Billing account not found. AccountId: %d", accountID))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found (001)."})
		return
	}

	type Invoice struct {
		Date          time.Time `json:"date"`
		Amount        float64   `json:"amount"`
		Transaction   string    `json:"transaction"`
		PaymentMethod string    `json:"payment_method"`
		InvoiceURL    string    `json:"invoice_url"`
	}

	invoices := []Invoice{}

	// Add trial period
	invoices = append(invoices, Invoice{
		Date:          billing.CreatedAt,
		Amount:        0,
		Transaction:   "Trial Period " + billing.CreatedAt.Format("1/2/06") + " - " + billing.TrialExpire.Format("1/2/06"),
		PaymentMethod: "",
		InvoiceURL:    "",
	})

	// If we have stripe key do this.
	if len(os.Getenv("STRIPE_SECRET_KEY")) > 0 {
		// Get charges by customer
		charges, err := stripe.GetChargesByCustomer(billing.StripeCustomer)

		if err != nil {
			services.Critical(fmt.Errorf("GetChargesByCustomer: Billing account not found. AccountId: %d", accountID))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found (001)."})
			return
		}

		// Loop through and add invoices
		for _, row := range charges {

			// Build invoice object
			tmp := Invoice{
				Date:          time.Unix(row.Created, 0),
				Amount:        float64(row.Amount / 100),
				Transaction:   "Charge",
				PaymentMethod: string(row.Source.Card.Brand) + " ending " + string(row.Source.Card.Last4),
				InvoiceURL:    "",
			}

			// is this a refund?
			if row.Refunded {
				tmp.Amount = float64(row.AmountRefunded/100) * -1
				tmp.Transaction = tmp.Transaction + " Refund"
			}

			// Add invoice information
			if row.Invoice != nil {
				inv, err := stripe.GetInvoice(row.Invoice.ID)

				if err == nil {
					tmp.InvoiceURL = inv.InvoicePDF

					// Get the date range
					start := time.Unix(inv.Lines.Data[0].Period.Start, 0).Format("1/2/06")
					end := time.Unix(inv.Lines.Data[0].Period.End, 0).Format("1/2/06")
					tmp.Transaction = "Subscription " + start + " - " + end
				}

			}

			invoices = append(invoices, tmp)
		}

	}

	// Return happy JSON
	c.JSON(200, invoices)
}

/* End File */
