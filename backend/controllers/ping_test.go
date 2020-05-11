//
//
// Date: 2020-05-11
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
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
// TestPingFromClient01 test ping from client
//
func TestPingFromClient01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Setup test data
	user := test.GetRandomUser(33)
	db.Save(&user)

	billing1 := test.GetRandomBilling(5, 33)
	billing1.StripeCustomer = ""
	db.Save(&billing1)
	account1 := test.GetRandomAccount(33)
	account1.OwnerId = user.Id
	account1.BillingId = 5
	db.Save(&account1)
	db.Save(&models.AcctToUsers{AccountId: account1.Id, UserId: user.Id})

	account2 := test.GetRandomAccount(34)
	account2.OwnerId = user.Id
	db.Save(&account2)
	db.Save(&models.AcctToUsers{AccountId: account2.Id, UserId: user.Id})

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/account/ping", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.GET("/api/v3/33/account/ping", c.PingFromClient)
	r.ServeHTTP(w, req)

	// Deconstruct the ping response
	type ping struct {
		Status string
	}

	// Grab result and convert to strut
	result := ping{}
	json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results.
	st.Expect(t, w.Code, 200)
	st.Expect(t, result.Status, "ok")
}

//
// TestPingFromClient02 test ping from client
//
func TestPingFromClient02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Setup test data
	user := test.GetRandomUser(33)
	db.Save(&user)

	billing1 := test.GetRandomBilling(5, 33)
	billing1.StripeCustomer = ""
	billing1.Status = "Expired"
	db.Save(&billing1)
	account1 := test.GetRandomAccount(33)
	account1.OwnerId = user.Id
	account1.BillingId = 5
	db.Save(&account1)
	db.Save(&models.AcctToUsers{AccountId: account1.Id, UserId: user.Id})

	account2 := test.GetRandomAccount(34)
	account2.OwnerId = user.Id
	db.Save(&account2)
	db.Save(&models.AcctToUsers{AccountId: account2.Id, UserId: user.Id})

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/account/ping", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.GET("/api/v3/33/account/ping", c.PingFromClient)
	r.ServeHTTP(w, req)

	// Deconstruct the ping response
	type ping struct {
		Status string
	}

	// Grab result and convert to strut
	result := ping{}
	json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results.
	st.Expect(t, w.Code, 200)
	st.Expect(t, result.Status, "expired")
}

//
// TestPingFromClient03 test ping from client
//
func TestPingFromClient03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Setup test data
	user := test.GetRandomUser(33)
	db.Save(&user)

	billing1 := test.GetRandomBilling(5, 33)
	billing1.StripeCustomer = ""
	billing1.Status = "Delinquent"
	db.Save(&billing1)
	account1 := test.GetRandomAccount(33)
	account1.OwnerId = user.Id
	account1.BillingId = 5
	db.Save(&account1)
	db.Save(&models.AcctToUsers{AccountId: account1.Id, UserId: user.Id})

	account2 := test.GetRandomAccount(34)
	account2.OwnerId = user.Id
	db.Save(&account2)
	db.Save(&models.AcctToUsers{AccountId: account2.Id, UserId: user.Id})

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/account/ping", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.GET("/api/v3/33/account/ping", c.PingFromClient)
	r.ServeHTTP(w, req)

	// Deconstruct the ping response
	type ping struct {
		Status string
	}

	// Grab result and convert to strut
	result := ping{}
	json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results.
	st.Expect(t, w.Code, 200)
	st.Expect(t, result.Status, "delinquent")
}

/* End File */
