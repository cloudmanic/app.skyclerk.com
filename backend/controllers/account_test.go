//
// Date: 2019-09-16
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	db, dbName, _ := models.NewTestDB("")
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

//
// TestUpdateAccount01 - update account
//
func TestUpdateAccount01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
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

	// Change account data.
	account1.Name = "Unit Test"
	account1.Currency = "BRL"
	account1.Locale = "pt-BR"

	// Get JSON
	putStr, _ := json.Marshal(account1)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/account", bytes.NewBuffer(putStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.PUT("/api/v3/33/account", c.UpdateAccount)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Account{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, result.Id, uint(33))
	st.Expect(t, result.OwnerId, uint(1))
	st.Expect(t, result.Name, "Unit Test")
	st.Expect(t, result.Locale, "pt-BR")
	st.Expect(t, result.Currency, "BRL")

	// Check database
	a := models.Account{}
	db.New().Find(&a, 33)

	// Test results.
	st.Expect(t, a.Id, uint(33))
	st.Expect(t, a.OwnerId, uint(1))
	st.Expect(t, a.Name, "Unit Test")
	st.Expect(t, a.Locale, "pt-BR")
	st.Expect(t, a.Currency, "BRL")
}

//
// TestUpdateAccount02 - update account - errors
//
func TestUpdateAccount02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
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

	// Change account data.
	account1.Name = ""
	account1.Currency = ""
	account1.Locale = ""

	// Get JSON
	putStr, _ := json.Marshal(account1)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/account", bytes.NewBuffer(putStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.PUT("/api/v3/33/account", c.UpdateAccount)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Account{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Body.String(), `{"errors":{"currency":"The currency field is required.","locale":"The locale field is required.","name":"The name field is required."}}`)
}

//
// TestUpdateAccount03 - Not Owner - errors
//
func TestUpdateAccount03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Setup test data
	user := test.GetRandomUser(33)
	db.Save(&user)

	account1 := test.GetRandomAccount(33)
	account1.OwnerId = uint(55)
	db.Save(&account1)
	db.Save(&models.AcctToUsers{AcctId: account1.Id, UserId: user.Id})

	// Change account data.
	account1.Name = ""
	account1.Currency = ""
	account1.Locale = ""

	// Get JSON
	putStr, _ := json.Marshal(account1)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/account", bytes.NewBuffer(putStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.PUT("/api/v3/33/account", c.UpdateAccount)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Body.String(), `{"error":"You must be the account owner."}`)
}

//
// TestUpdateAccount04 - Get updating account owner
//
func TestUpdateAccount04(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Setup test data
	user1 := test.GetRandomUser(33)
	user2 := test.GetRandomUser(33)
	user3 := test.GetRandomUser(33)
	db.Save(&user1)
	db.Save(&user2)
	db.Save(&user3)

	account1 := test.GetRandomAccount(33)
	account1.OwnerId = user1.Id
	db.Save(&account1)
	db.Save(&models.AcctToUsers{AcctId: account1.Id, UserId: user1.Id})
	db.Save(&models.AcctToUsers{AcctId: account1.Id, UserId: user2.Id})
	db.Save(&models.AcctToUsers{AcctId: account1.Id, UserId: user3.Id})

	// Change account data.
	account1.OwnerId = user3.Id

	// Get JSON
	putStr, _ := json.Marshal(account1)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/account", bytes.NewBuffer(putStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user1.Id))
	})
	r.PUT("/api/v3/33/account", c.UpdateAccount)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Account{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, result.Id, uint(33))
	st.Expect(t, result.OwnerId, uint(3))

	// Check database
	a := models.Account{}
	db.New().Find(&a, 33)

	// Test results.
	st.Expect(t, a.Id, uint(33))
	st.Expect(t, a.OwnerId, uint(3))
}

//
// TestUpdateAccount05 - Get updating account owner - error
//
func TestUpdateAccount05(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Setup test data
	user1 := test.GetRandomUser(33)
	user2 := test.GetRandomUser(33)
	user3 := test.GetRandomUser(33)
	db.Save(&user1)
	db.Save(&user2)
	db.Save(&user3)

	account1 := test.GetRandomAccount(33)
	account1.OwnerId = user1.Id
	db.Save(&account1)
	db.Save(&models.AcctToUsers{AcctId: account1.Id, UserId: user1.Id})
	db.Save(&models.AcctToUsers{AcctId: account1.Id, UserId: user2.Id})
	db.Save(&models.AcctToUsers{AcctId: uint(44), UserId: user3.Id})

	// Change account data.
	account1.OwnerId = user3.Id

	// Get JSON
	putStr, _ := json.Marshal(account1)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/account", bytes.NewBuffer(putStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user1.Id))
	})
	r.PUT("/api/v3/33/account", c.UpdateAccount)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Body.String(), `{"errors":{"owner_id":"Invalid owner_id was posted."}}`)
}

//
// TestUpdateAccount06 - Get updating account owner - error, more than one field.
//
func TestUpdateAccount06(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Setup test data
	user1 := test.GetRandomUser(33)
	user2 := test.GetRandomUser(33)
	user3 := test.GetRandomUser(33)
	db.Save(&user1)
	db.Save(&user2)
	db.Save(&user3)

	account1 := test.GetRandomAccount(33)
	account1.OwnerId = user1.Id
	account1.Name = ""
	db.Save(&account1)
	db.Save(&models.AcctToUsers{AcctId: account1.Id, UserId: user1.Id})
	db.Save(&models.AcctToUsers{AcctId: account1.Id, UserId: user2.Id})
	db.Save(&models.AcctToUsers{AcctId: uint(44), UserId: user3.Id})

	// Change account data.
	account1.OwnerId = user3.Id

	// Get JSON
	putStr, _ := json.Marshal(account1)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/account", bytes.NewBuffer(putStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user1.Id))
	})
	r.PUT("/api/v3/33/account", c.UpdateAccount)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"errors":{"name":"The name field is required.","owner_id":"Invalid owner_id was posted."}}`)
}

//
// TestClearAccount01 - Clear account.
//
func TestClearAccount01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Setup test data
	user1 := test.GetRandomUser(33)
	user2 := test.GetRandomUser(33)
	user3 := test.GetRandomUser(33)
	db.Save(&user1)
	db.Save(&user2)
	db.Save(&user3)

	account1 := test.GetRandomAccount(33)
	account1.OwnerId = user1.Id
	db.Save(&account1)
	db.Save(&models.AcctToUsers{AcctId: account1.Id, UserId: user1.Id})
	db.Save(&models.AcctToUsers{AcctId: account1.Id, UserId: user2.Id})
	db.Save(&models.AcctToUsers{AcctId: uint(44), UserId: user3.Id})

	// Create like 10 ledger entries.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		db.LedgerCreate(&l)
	}

	// Create like 10 ledger entries. - Different account.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(34)
		db.LedgerCreate(&l)
	}

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/account/clear", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user1.Id))
	})
	r.POST("/api/v3/33/account/clear", c.ClearAccount)
	r.ServeHTTP(w, req)

	// Get the ledger entries. There should not be any with account 33
	l := []models.Ledger{}
	db.Find(&l)
	for _, row := range l {
		st.Expect(t, (row.AccountId == uint(33)), false)
	}

	// Get the Category entries.
	cats := []models.Category{}
	db.Where("CategoriesAccountId = ?", 33).Find(&cats)

	// Test results
	st.Expect(t, w.Code, 204)
	st.Expect(t, len(l), 10)
	st.Expect(t, len(cats), 23)
}

//
// TestClearAccount02 - Clear account - not owner.
//
func TestClearAccount02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Setup test data
	user1 := test.GetRandomUser(33)
	user2 := test.GetRandomUser(33)
	user3 := test.GetRandomUser(33)
	db.Save(&user1)
	db.Save(&user2)
	db.Save(&user3)

	account1 := test.GetRandomAccount(33)
	account1.OwnerId = user1.Id
	db.Save(&account1)
	db.Save(&models.AcctToUsers{AcctId: account1.Id, UserId: user1.Id})
	db.Save(&models.AcctToUsers{AcctId: account1.Id, UserId: user2.Id})
	db.Save(&models.AcctToUsers{AcctId: uint(44), UserId: user3.Id})

	// Create like 5 ledger entries.
	for i := 0; i < 5; i++ {
		l := test.GetRandomLedger(33)
		db.LedgerCreate(&l)
	}

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/account/clear", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user2.Id))
	})
	r.POST("/api/v3/33/account/clear", c.ClearAccount)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"You must be the account owner."}`)

	// Get the ledger entries. There should not be any with account 33
	l := []models.Ledger{}
	db.Find(&l)

	// Test results
	st.Expect(t, len(l), 5)
}

//
// TestDeleteAccount01 - Delete account.
//
func TestDeleteAccount01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Setup test data
	user1 := test.GetRandomUser(33)
	user2 := test.GetRandomUser(33)
	user3 := test.GetRandomUser(33)
	user4 := test.GetRandomUser(33)
	db.Save(&user1)
	db.Save(&user2)
	db.Save(&user3)
	db.Save(&user4)

	account1 := test.GetRandomAccount(33)
	account1.OwnerId = user1.Id
	db.Save(&account1)
	db.Save(&models.AcctToUsers{AcctId: account1.Id, UserId: user1.Id})
	db.Save(&models.AcctToUsers{AcctId: account1.Id, UserId: user2.Id})
	db.Save(&models.AcctToUsers{AcctId: account1.Id, UserId: user4.Id})

	account2 := test.GetRandomAccount(34)
	account2.OwnerId = user1.Id
	db.Save(&account2)
	db.Save(&models.AcctToUsers{AcctId: account2.Id, UserId: user1.Id})
	db.Save(&models.AcctToUsers{AcctId: account2.Id, UserId: user2.Id})

	// Create like 10 ledger entries.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		db.LedgerCreate(&l)
	}

	// Create like 10 ledger entries. - Different account.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(34)
		db.LedgerCreate(&l)
	}

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/account/delete", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user1.Id))
	})
	r.POST("/api/v3/33/account/delete", c.DeleteAccount)
	r.ServeHTTP(w, req)

	// Get the ledger entries. There should not be any with account 33
	l := []models.Ledger{}
	db.Find(&l)
	for _, row := range l {
		st.Expect(t, (row.AccountId == uint(33)), false)
	}

	// Get the Category entries.
	cats := []models.Category{}
	db.Where("CategoriesAccountId = ?", 33).Find(&cats)

	// Get the AcctToUsers entries.
	a2u := []models.AcctToUsers{}
	db.Where("acct_id = ?", 33).Find(&a2u)

	// Get the Account entries.
	acc := []models.Account{}
	db.Where("id = ?", 33).Find(&acc)

	// Grab result and convert to strut
	results := []models.Account{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, len(l), 10)
	st.Expect(t, len(a2u), 0)
	st.Expect(t, len(acc), 0)
	st.Expect(t, len(cats), 0)
	st.Expect(t, len(results), 1)
}

//
// TestDeleteAccount02 - Delete account.
//
func TestDeleteAccount02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Setup test data
	user1 := test.GetRandomUser(33)
	db.Save(&user1)

	account1 := test.GetRandomAccount(33)
	account1.OwnerId = user1.Id
	db.Save(&account1)
	db.Save(&models.AcctToUsers{AcctId: account1.Id, UserId: user1.Id})

	// Create like 10 ledger entries.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		db.LedgerCreate(&l)
	}

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/account/delete", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user1.Id))
	})
	r.POST("/api/v3/33/account/delete", c.DeleteAccount)
	r.ServeHTTP(w, req)

	// Get the Category entries.
	cats := []models.Category{}
	db.Where("CategoriesAccountId = ?", 33).Find(&cats)

	// Get the AcctToUsers entries.
	a2u := []models.AcctToUsers{}
	db.Where("acct_id = ?", 33).Find(&a2u)

	// Get the Account entries.
	acc := []models.Account{}
	db.Where("id = ?", 33).Find(&acc)

	// Grab result and convert to strut
	results := []models.Account{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, len(a2u), 0)
	st.Expect(t, len(acc), 0)
	st.Expect(t, len(cats), 0)
	st.Expect(t, len(results), 0)
}

//
// TestNewAccount01 - Add new account.
//
func TestNewAccount01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Setup test data
	user1 := test.GetRandomUser(33)
	db.Save(&user1)

	account1 := test.GetRandomAccount(33)
	account1.OwnerId = user1.Id
	db.Save(&account1)
	db.Save(&models.AcctToUsers{AcctId: account1.Id, UserId: user1.Id})

	// Add billing profile
	b := models.Billing{}
	db.Save(&b)

	// Add account to billing
	atb := models.AcctToBilling{AcctId: account1.Id, BillingId: b.Id}
	db.Save(&atb)

	// Get JSON
	postStr := fmt.Sprintf(`{ "name": "%s" }`, "Unit Test Account")

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/account/new", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user1.Id))
	})
	r.POST("/api/v3/33/account/new", c.NewAccount)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Account{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Get the Account entries.
	au := []models.AcctToUsers{}
	db.Where("acct_id = ?", 34).Find(&au)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, result.Name, "Unit Test Account")
	st.Expect(t, result.Id, uint(34))
	st.Expect(t, len(au), 1)
	st.Expect(t, au[0].UserId, uint(1))
}

/* End File */