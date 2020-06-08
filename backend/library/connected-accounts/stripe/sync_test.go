//
// Date: 2020-06-07
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

package stripe

import (
	"testing"

	"app.skyclerk.com/backend/library/test"
	"app.skyclerk.com/backend/models"
)

//
// TestSync01 will sync stripe
//
func TestSync01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("testing_db")
	defer models.TestingTearDown(db, dbName)

	// Setup test data
	user := test.GetRandomUser(33)
	db.Save(&user)

	account := test.GetRandomAccount(33)
	account.OwnerId = user.Id
	db.Save(&account)
	db.Save(&models.AcctToUsers{AccountId: account.Id, UserId: user.Id})

	// Connected account model.
	ac := models.ConnectedAccounts{
		AccountID:            33,
		Connection:           "Stripe",
		StripeUserID:         "Ris6D3eYxtXbIkRz5K7aJBLiGkTOvuuD", // Cloudmanic Labs stripe test
		StripeScope:          "read_only",
		StripeLastItem:       1487731749, // Sample last item.
		StripePublishableKey: "pk_test_1Ris6D3eYxtXbIkRz5K7aJBLiGkTOvuuDzkgPVNu7mF1aON92g9n6xREVSPcF134HGoKuuh9sCwgt8Ai1D6ApFaR100JiHDosXn",
	}
	db.New().Save(&ac)

	// Do sync
	Sync(db, account, ac)
}

/* End File */
