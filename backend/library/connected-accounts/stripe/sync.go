//
// Date: 2020-06-07
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

package stripe

import (
	"fmt"
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
func Sync(db models.Datastore, connectedAccount models.ConnectedAccounts) {
	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		services.LogInfo("No STRIPE_SECRET_KEY found in stripe.sync.")
		return
	}

	// Track last item
	lastItem := int64(0)
	processCount := 0

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

	// Account we are trying to access. (XXX is Cloudmanic account)
	if connectedAccount.StripeUserID != "XXXX" {
		params.SetStripeAccount(connectedAccount.StripeUserID)
	}

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

		// XXX is Cloudmanic account
		if connectedAccount.StripeUserID != "XXXX" {
			p.SetStripeAccount(connectedAccount.StripeUserID)
		}

		bt, _ := balancetransaction.Get(c.BalanceTransaction.ID, p)

		// Get the customer
		y := &stripe.CustomerParams{}

		// XXX is Cloudmanic account
		if connectedAccount.StripeUserID != "XXXX" {
			y.SetStripeAccount(connectedAccount.StripeUserID)
		}

		cust, _ := customer.Get(c.Customer.ID, y)

		// Process this transaction
		processTransaction(db, connectedAccount, c.ID, c.Amount, bt.Fee, c.Created, c.Customer.ID, cust.Email, cust.Name, cust.Description)

		// Flag the last item.
		if lastItem < c.Created {
			lastItem = c.Created
		}

		// Counter
		processCount++
	}

	// Update sync info.
	connectedAccount.StripeLastSync = time.Now()

	// Only set this when we have new items
	if lastItem > 0 {
		connectedAccount.StripeLastItem = lastItem
	}

	// Update connected app database entry.
	db.New().Save(&connectedAccount)

	// Logging
	services.InfoMsg(fmt.Sprintf("Processed %d Stripe items for Account: %d, Last Item: %d", processCount, connectedAccount.AccountID, lastItem))
}

//
// processTransaction will store the transaction in the database.
//
func processTransaction(
	db models.Datastore,
	connectedAccount models.ConnectedAccounts,
	tranID string,
	amount int64,
	fee int64,
	createdAt int64,
	custID string,
	custEmail string,
	custName string,
	custDesc string) {
	// Do a quick check to make sure we do not already have this ledger entry.
	lg := models.Ledger{}
	db.New().Where("LedgerAccountId = ? AND LedgerStripeId = ?", connectedAccount.AccountID, tranID).First(&lg)

	if lg.Id > 0 {
		services.LogInfo("Stripe entry already in ledger: " + tranID)
		return
	}

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
	db.New().Where("ContactsAccountId = ? AND stripe_cust_id = ?", connectedAccount.AccountID, custID).First(&contact)

	// Update record and save.
	if contact.Id == 0 {
		contact.AccountId = connectedAccount.AccountID
		contact.Name = name
		contact.Email = custEmail
		contact.StripeCustID = custID
		db.New().Save(&contact)
	}

	// Setup the new contact for Stripe
	feeContact := models.Contact{}

	// See if we have this record.
	db.New().Where("ContactsAccountId = ? AND ContactsName = ?", connectedAccount.AccountID, "Stripe").First(&feeContact)

	// Update record and save.
	if feeContact.Id == 0 {
		feeContact.AccountId = connectedAccount.AccountID
		feeContact.Name = "Stripe"
		feeContact.Website = "https://stripe.com"
		db.New().Save(&feeContact)
	}

	// Get income category TODO(spicer): Let the user select this category
	incomeCat := db.GetOrCreateCategory(connectedAccount.AccountID, "Stripe Charges", "2")

	// Get fee category
	feeCat := db.GetOrCreateCategory(connectedAccount.AccountID, "Payment Processing Fee", "1")

	// Create or get label
	label := db.GetOrCreateLabel(connectedAccount.AccountID, "stripe")

	// Add the ledger entry to the database
	ledger := models.Ledger{
		AccountId:  connectedAccount.AccountID,
		ContactId:  contact.Id,
		Date:       time.Unix(createdAt, 0),
		Amount:     float64(float64(amount) / float64(100)),
		CategoryId: incomeCat.Id,
		StripeId:   tranID,
		Note:       "Stripe Import of charge - " + tranID,
		Labels:     []models.Label{label},
	}
	db.New().Save(&ledger)

	// Insert the stripe fee.
	feeObj := models.Ledger{
		AccountId:  connectedAccount.AccountID,
		ContactId:  feeContact.Id,
		Date:       time.Unix(createdAt, 0),
		Amount:     float64(float64(fee)/float64(100)) * -1,
		CategoryId: feeCat.Id,
		StripeId:   tranID,
		Note:       "Stripe Fee of charge - " + tranID,
		Labels:     []models.Label{label},
	}
	db.New().Save(&feeObj)

}

/* End File */
