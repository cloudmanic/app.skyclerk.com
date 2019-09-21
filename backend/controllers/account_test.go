//
// Date: 2019-09-16
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"app.skyclerk.com/backend/library/test"
	"app.skyclerk.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// TestGetAccount01 - test account
//
func TestGetAccount01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("testing_db")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Setup test data
	user := test.GetRandomUser(33)
	db.Save(&user)

	account1 := test.GetRandomAccount(33)
	account1.OwnerId = user.Id
	db.Save(&account1)
	db.Save(&models.AcctToUsers{AcctId: account1.Id, UserId: user.Id})

	account2 := test.GetRandomAccount(34)
	account2.OwnerId = user.Id
	db.Save(&account2)
	db.Save(&models.AcctToUsers{AcctId: account2.Id, UserId: user.Id})

	account3 := test.GetRandomAccount(105)
	account3.OwnerId = user.Id
	db.Save(&account3)
	db.Save(&models.AcctToUsers{AcctId: account3.Id, UserId: user.Id})

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/account", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/33/account", c.GetAccount)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Account{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, result.Id, uint(33))
	st.Expect(t, result.OwnerId, uint(1))
	st.Expect(t, result.Name, account1.Name)
	st.Expect(t, result.Locale, "en-US")
	st.Expect(t, result.Currency, "USD")
}

/* End File */
