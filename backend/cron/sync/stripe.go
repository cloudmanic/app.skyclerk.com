//
// Date: 2020-06-07
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

package sync

import (
	"app.skyclerk.com/backend/library/connected-accounts/stripe"
	"app.skyclerk.com/backend/models"
	"app.skyclerk.com/backend/services"
)

//
// StripeSync will sync in transactions from stripe.
//
func StripeSync(db models.Datastore) {
	services.InfoMsg("Starting stripe entry sync.")

	// Get all connected accounts
	ca := []models.ConnectedAccounts{}
	db.New().Where("connection = ?", "Stripe").Find(&ca)

	// Loop through all the stripe connections.
	for _, row := range ca {
		// Do sync
		stripe.Sync(db, row)
	}

	services.InfoMsg("All stripe entries were imported.")
}

/* End File */
