//
// Date: 2020-06-07
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

package stripe

import (
	"os"
	"time"

	"app.skyclerk.com/backend/models"
	"app.skyclerk.com/backend/services"
	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/balancetransaction"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/customer"
)

//
// Sync will bring in all transactions for a stripe account.
//
func Sync(db models.Datastore, account models.Account, connectedAccount models.ConnectedAccounts) {
	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		services.LogInfo("No STRIPE_SECRET_KEY found in stripe.sync.")
		return
	}

	// Set stripe key
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Turn verbose logging off
	stripe.DefaultLeveledLogger = &stripe.LeveledLogger{
		Level: stripe.LevelError,
	}

	// Set parms
	params := &stripe.ChargeListParams{
		CreatedRange: &stripe.RangeQueryParams{
			GreaterThan: connectedAccount.StripeLastItem,
		},
	}

	// Account we are trying to access.
	params.SetStripeAccount(connectedAccount.StripeUserID)

	// Filter
	params.Filters.AddFilter("limit", "", "100")

	// This will bering in all charges - auto paging.
	i := charge.List(params)
	for i.Next() {
		c := i.Charge()

		// Only collect succeeded charges
		if c.Status != "succeeded" {
			continue
		}

		// Get the balance transaction
		p := &stripe.BalanceTransactionParams{}
		p.SetStripeAccount(connectedAccount.StripeUserID)
		bt, _ := balancetransaction.Get(c.BalanceTransaction.ID, p)

		// Get the customer
		y := &stripe.CustomerParams{}
		y.SetStripeAccount(connectedAccount.StripeUserID)
		cust, _ := customer.Get(c.Customer.ID, y)

		// Process this transaction
		processTransaction(db, account, c.ID, c.Amount, bt.Fee, c.Created, c.Customer.ID, cust.Email, cust.Name, cust.Description)
	}

}

//
// processTransaction will store the transaction in the database.
//
func processTransaction(
	db models.Datastore,
	account models.Account,
	tranID string,
	amount int64,
	fee int64,
	createdAt int64,
	custID string,
	custEmail string,
	custName string,
	custDesc string) {
	// Figure out the customer name
	name := ""

	if len(custName) > 0 {
		name = custName
	} else {
		name = "Stripe Customer - " + custID
	}

	// Setup the new contact.
	contact := models.Contact{}

	// See if we have this record.
	db.New().Where("ContactsAccountId = ? AND stripe_cust_id = ?", account.Id, custID).First(&contact)

	// Update record and save.
	contact.AccountId = account.Id
	contact.Name = name
	contact.Email = custEmail
	contact.StripeCustID = custID
	db.New().Save(&contact)

	// Add the ledger entry to the database
	ledger := models.Ledger{
		AccountId:  account.Id,
		ContactId:  contact.Id,
		Date:       time.Unix(createdAt, 0),
		Amount:     float64(float64(amount) / float64(100)),
		CategoryId: 0,
		StripeId:   tranID,
		Note:       "Stripe Import of charge - " + tranID,
	}
	db.New().Save(&ledger)

	// Insert the stripe fee.
	feeObj := models.Ledger{
		AccountId:  account.Id,
		ContactId:  contact.Id, // Make this the stripe contact.
		Date:       time.Unix(createdAt, 0),
		Amount:     float64(float64(fee)/float64(100)) * -1,
		CategoryId: 0,
		StripeId:   tranID,
		Note:       "Stripe Fee of charge - " + tranID,
	}
	db.New().Save(&feeObj)

}

/* End File */
