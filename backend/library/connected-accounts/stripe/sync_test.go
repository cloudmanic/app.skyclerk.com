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
	Sync(db, ac)

	// Verify the import
	l := []models.Ledger{}
	db.New().Preload("Contact").Preload("Category").Preload("Labels").Where("LedgerAccountId = ?", 33).Find(&l)
	st.Expect(t, l[36].Id, uint(37))
	st.Expect(t, l[36].StripeId, "ch_AJuAINDXfQOdVA")
	st.Expect(t, l[36].Note, "Stripe Import of charge - ch_AJuAINDXfQOdVA")
	st.Expect(t, l[36].Amount, 15.00)
	st.Expect(t, l[36].CategoryId, uint(1))
	st.Expect(t, l[36].Contact.Name, "Stripe Customer - cus_4zqHAQ3D03mPSf")
	st.Expect(t, l[36].Labels[0].Name, "stripe")
	st.Expect(t, l[37].Id, uint(38))
	st.Expect(t, l[37].StripeId, "ch_AJuAINDXfQOdVA")
	st.Expect(t, l[37].Note, "Stripe Fee of charge - ch_AJuAINDXfQOdVA")
	st.Expect(t, l[37].Amount, -0.74)
	st.Expect(t, l[37].CategoryId, uint(2))
	st.Expect(t, l[37].Contact.Name, "Stripe")
	st.Expect(t, l[37].Contact.Website, "https://stripe.com")
	st.Expect(t, l[37].Labels[0].Name, "stripe")
}

//
// TestSync02 will sync stripe - but our compay
//
func TestSync02(t *testing.T) {
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
		StripeUserID:         "XXXX", // Use to mark this as our api key
		StripeScope:          "read_only",
		StripeLastItem:       1591634741, // Sample last item.
		StripePublishableKey: "",
	}
	db.New().Save(&ac)

	// Do sync
	Sync(db, ac)

	// Verify the import
	l := []models.Ledger{}
	db.New().Preload("Contact").Preload("Category").Preload("Labels").Where("LedgerAccountId = ?", 33).Order("LedgerId ASC").Find(&l)
	st.Expect(t, l[61].Id, uint(62))
	st.Expect(t, l[61].StripeId, "ch_1Gro6XKicNxxeq0znFuSNHmO")
	st.Expect(t, l[61].Note, "Stripe Fee of charge - ch_1Gro6XKicNxxeq0znFuSNHmO")
	st.Expect(t, l[61].Amount, -0.47)
	st.Expect(t, l[61].CategoryId, uint(2))
	st.Expect(t, l[61].Contact.Name, "Stripe")
	st.Expect(t, l[61].Labels[0].Name, "stripe")
	st.Expect(t, l[60].Id, uint(61))
	st.Expect(t, l[60].StripeId, "ch_1Gro6XKicNxxeq0znFuSNHmO")
	st.Expect(t, l[60].Note, "Stripe Import of charge - ch_1Gro6XKicNxxeq0znFuSNHmO")
	st.Expect(t, l[60].Amount, 6.00)
	st.Expect(t, l[60].CategoryId, uint(1))
	st.Expect(t, l[60].Contact.Name, "Stripe Customer - cus_HQfRLXUsmvqoKY")
	st.Expect(t, l[60].Labels[0].Name, "stripe")
}

/* End File */
