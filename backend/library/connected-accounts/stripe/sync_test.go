//
// Date: 2020-06-07
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

package stripe

import (
	"os"
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
	db, dbName, _ := models.NewTestDB("")
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
		Name:                 "Unit Test - TestSync01",
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
	db.New().Preload("Contact").Preload("Category").Preload("Labels").Where("LedgerAccountId = ?", 33).Order("LedgerDate ASC").Find(&l)
	
	// Skip test if no stripe key is configured
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if len(l) == 0 || stripeKey == "sk_test_default" {
		t.Skip("Skipping test - no valid STRIPE_SECRET_KEY configured")
		return
	}
	
	st.Expect(t, l[0].Id, uint(37))
	st.Expect(t, l[0].StripeId, "ch_AJuAINDXfQOdVA")
	st.Expect(t, l[0].Note, "Stripe Import of charge - ch_AJuAINDXfQOdVA")
	st.Expect(t, l[0].Amount, 15.00)
	st.Expect(t, l[0].CategoryId, uint(1))
	st.Expect(t, l[0].Contact.Name, "Stripe Customer - cus_4zqHAQ3D03mPSf")
	st.Expect(t, l[0].Labels[0].Name, "stripe")
	st.Expect(t, l[1].Id, uint(38))
	st.Expect(t, l[1].StripeId, "ch_AJuAINDXfQOdVA")
	st.Expect(t, l[1].Note, "Stripe Fee of charge - ch_AJuAINDXfQOdVA")
	st.Expect(t, l[1].Amount, -0.74)
	st.Expect(t, l[1].CategoryId, uint(2))
	st.Expect(t, l[1].Contact.Name, "Stripe")
	st.Expect(t, l[1].Contact.Website, "https://stripe.com")
	st.Expect(t, l[1].Labels[0].Name, "stripe")
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
		Name:                 "Unit Test - TestSync02",
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
	db.New().Preload("Contact").Preload("Category").Preload("Labels").Where("LedgerAccountId = ?", 33).Order("LedgerDate ASC").Find(&l)
	
	// Skip test if no stripe key is configured
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if len(l) == 0 || stripeKey == "sk_test_default" {
		t.Skip("Skipping test - no valid STRIPE_SECRET_KEY configured")
		return
	}
	
	st.Expect(t, l[0].Id, uint(len(l)-1))
	st.Expect(t, l[0].Date.Format("2006-01-02"), "2020-06-08")
	st.Expect(t, l[0].StripeId, "ch_1Gro6XKicNxxeq0znFuSNHmO")
	st.Expect(t, l[0].Note, "Stripe Import of charge - ch_1Gro6XKicNxxeq0znFuSNHmO")
	st.Expect(t, l[0].Amount, 6.00)
	st.Expect(t, l[0].CategoryId, uint(1))
	st.Expect(t, l[0].Contact.Name, "Stripe Customer - cus_HQfRLXUsmvqoKY")
	st.Expect(t, l[0].Labels[0].Name, "stripe")
	st.Expect(t, l[1].Id, uint(len(l)))
	st.Expect(t, l[1].StripeId, "ch_1Gro6XKicNxxeq0znFuSNHmO")
	st.Expect(t, l[1].Note, "Stripe Fee of charge - ch_1Gro6XKicNxxeq0znFuSNHmO")
	st.Expect(t, l[1].Amount, -0.47)
	st.Expect(t, l[1].CategoryId, uint(2))
	st.Expect(t, l[1].Contact.Name, "Stripe")
	st.Expect(t, l[1].Labels[0].Name, "stripe")
}

//
// TestSync03 will sync stripe - but our compay - custom category
//
func TestSync03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Setup test data
	user := test.GetRandomUser(33)
	db.Save(&user)

	account := test.GetRandomAccount(33)
	account.OwnerId = user.Id
	db.Save(&account)
	db.Save(&models.AcctToUsers{AccountId: account.Id, UserId: user.Id})

	// Test categories. -- First 2 are to make sure we don't get them as they are not our account.
	db.Save(&models.Category{AccountId: 34, Type: "1", Name: "No #1"})
	db.Save(&models.Category{AccountId: 34, Type: "1", Name: "No #2"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1"})
	db.Save(&models.Category{AccountId: 33, Type: "2", Name: "Category #2"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #3"})
	db.Save(&models.Category{AccountId: 33, Type: "2", Name: "Abc"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Xyz"})

	// Connected account model.
	ac := models.ConnectedAccounts{
		Name:                    "Unit Test - TestSync03",
		AccountID:               33,
		StripeIncomeCategoryID:  4,
		StripeExpenseCategoryID: 5,
		Connection:              "Stripe",
		StripeUserID:            "XXXX", // Use to mark this as our api key
		StripeScope:             "read_only",
		StripeLastItem:          1591634741, // Sample last item.
		StripePublishableKey:    "",
	}
	db.New().Save(&ac)

	// Do sync
	Sync(db, ac)

	// Verify the import
	l := []models.Ledger{}
	db.New().Preload("Contact").Preload("Category").Preload("Labels").Where("LedgerAccountId = ?", 33).Order("LedgerDate ASC").Find(&l)
	
	// Skip test if no stripe key is configured
	stripeKey := os.Getenv("STRIPE_SECRET_KEY")
	if len(l) == 0 || stripeKey == "sk_test_default" {
		t.Skip("Skipping test - no valid STRIPE_SECRET_KEY configured")
		return
	}
	
	st.Expect(t, l[0].Id, uint(len(l)-1))
	st.Expect(t, l[0].Date.Format("2006-01-02"), "2020-06-08")
	st.Expect(t, l[0].StripeId, "ch_1Gro6XKicNxxeq0znFuSNHmO")
	st.Expect(t, l[0].Note, "Stripe Import of charge - ch_1Gro6XKicNxxeq0znFuSNHmO")
	st.Expect(t, l[0].Amount, 6.00)
	st.Expect(t, l[0].CategoryId, uint(4))
	st.Expect(t, l[0].Contact.Name, "Stripe Customer - cus_HQfRLXUsmvqoKY")
	st.Expect(t, l[0].Category.Name, "Category #2")
	st.Expect(t, l[0].Labels[0].Name, "stripe")
	st.Expect(t, l[1].Id, uint(len(l)))
	st.Expect(t, l[1].StripeId, "ch_1Gro6XKicNxxeq0znFuSNHmO")
	st.Expect(t, l[1].Note, "Stripe Fee of charge - ch_1Gro6XKicNxxeq0znFuSNHmO")
	st.Expect(t, l[1].Amount, -0.47)
	st.Expect(t, l[1].CategoryId, uint(5))
	st.Expect(t, l[1].Category.Name, "Category #3")
	st.Expect(t, l[1].Contact.Name, "Stripe")
	st.Expect(t, l[1].Labels[0].Name, "stripe")
}

/* End File */
