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
	"github.com/nbio/st"
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

	// Verify the import
	l := []models.Ledger{}
	db.New().Preload("Contact").Preload("Category").Preload("Labels").Where("LedgerAccountId = ?", 33).Find(&l)
	st.Expect(t, l[0].Id, uint(1))
	st.Expect(t, l[0].StripeId, "ch_FRYSWXArFoJY05")
	st.Expect(t, l[0].Note, "Stripe Import of charge - ch_FRYSWXArFoJY05")
	st.Expect(t, l[0].Amount, 9012.00)
	st.Expect(t, l[0].CategoryId, uint(1))
	st.Expect(t, l[0].Contact.Name, "Blah Matthews")
	st.Expect(t, l[0].Labels[0].Name, "stripe")
	st.Expect(t, l[1].Id, uint(2))
	st.Expect(t, l[1].StripeId, "ch_FRYSWXArFoJY05")
	st.Expect(t, l[1].Note, "Stripe Fee of charge - ch_FRYSWXArFoJY05")
	st.Expect(t, l[1].Amount, -261.65)
	st.Expect(t, l[1].CategoryId, uint(2))
	st.Expect(t, l[1].Contact.Name, "stripe")
	st.Expect(t, l[1].Contact.Website, "https://stripe.com")
	st.Expect(t, l[1].Labels[0].Name, "stripe")

}

/* End File */
