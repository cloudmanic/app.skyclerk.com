//
// Date: 11/3/2018
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package admin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"app.skyclerk.com/backend/library/test"
	"app.skyclerk.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// TestGetAccounts01 - return a list of accounts with meta data.
//
func TestGetAccounts01(t *testing.T) {
	// Data map
	dMap := make(map[uint]models.Ledger)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("testing_db")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Setup test data
	user1 := test.GetRandomUser(33)
	db.Save(&user1)
	user2 := test.GetRandomUser(34)
	db.Save(&user2)
	user3 := test.GetRandomUser(105)
	db.Save(&user3)

	account1 := test.GetRandomAccount(33)
	account1.OwnerId = user1.Id
	account1.BillingId = 5
	account1.LastActivity = time.Now().Add(time.Hour * 1)
	db.Save(&account1)
	db.Save(&models.AcctToUsers{AccountId: account1.Id, UserId: user1.Id})
	billing1 := test.GetRandomBilling(5, 33)
	db.Save(&billing1)

	account2 := test.GetRandomAccount(34)
	account2.OwnerId = user2.Id
	account2.BillingId = 6
	account2.LastActivity = time.Now().Add(time.Hour * 2)
	db.Save(&account2)
	db.Save(&models.AcctToUsers{AccountId: account2.Id, UserId: user2.Id})
	billing2 := test.GetRandomBilling(6, 34)
	db.Save(&billing2)

	account3 := test.GetRandomAccount(105)
	account3.OwnerId = user3.Id
	account3.BillingId = 7
	account3.LastActivity = time.Now().Add(time.Hour * 3)
	db.Save(&account3)
	db.Save(&models.AcctToUsers{AccountId: account3.Id, UserId: user3.Id})
	billing3 := test.GetRandomBilling(7, 105)
	db.Save(&billing3)

	// Create like 10 ledger entries. Diffent account.
	for i := 0; i < 11; i++ {
		l := test.GetRandomLedger(33)
		db.LedgerCreate(&l)
	}

	// Create like 105 ledger entries.
	for i := 0; i < 105; i++ {
		l := test.GetRandomLedger(34)
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Create like 10 ledger entries. Diffent account.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(105)
		db.LedgerCreate(&l)
	}

	// Setup request
	req, _ := http.NewRequest("GET", "/api/admin/accounts", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/admin/accounts", c.GetAccounts)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := []Account{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, result[0].AccountID, uint(105))
	st.Expect(t, result[1].AccountID, uint(34))
	st.Expect(t, result[2].AccountID, uint(33))
	st.Expect(t, result[0].LedgerCount, 10)
	st.Expect(t, result[1].LedgerCount, 105)
	st.Expect(t, result[2].LedgerCount, 11)
	st.Expect(t, result[0].Email, user3.Email)
	st.Expect(t, result[1].Email, user2.Email)
	st.Expect(t, result[2].Email, user1.Email)
}

/* End File */
