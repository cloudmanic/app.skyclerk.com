//
// Date: 2019-06-26
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
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
// Test get Activities 01
//
func TestGetActivities01(t *testing.T) {
	// Data map
	dMap := make(map[uint]models.Activity)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test user
	user := test.GetRandomUser(33)
	db.Save(&user)

	// Create like 105 ledger entries. This will create Activities
	for i := 0; i < 105; i++ {
		l := test.GetRandomLedger(33)
		db.LedgerCreate(&l)

		// Set the ledger type
		ledgerType := "expense"

		if l.Amount > 0 {
			ledgerType = "income"
		}

		// Get the contact name.
		contactName := l.Contact.Name

		if len(contactName) == 0 {
			contactName = l.Contact.FirstName + " " + l.Contact.LastName
		}

		// Add to the activity log
		y := models.Activity{
			AccountId: uint(33),
			UserId:    user.Id,
			Action:    ledgerType,
			SubAction: "create",
			Name:      contactName,
			Amount:    l.Amount,
			LedgerId:  l.Id,
		}

		db.Create(&y)

		dMap[l.Id] = y
	}

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/activities", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.GET("/api/v3/:account/activities", c.GetActivities)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []models.Activity{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, len(results), 100)
	st.Expect(t, results[0].Id, uint(105))
	st.Expect(t, results[1].Id, uint(104))
	st.Expect(t, results[2].Id, uint(103))
	st.Expect(t, results[0].Message, fmt.Sprintf("%s %s an %s ledger entry of %.2f from %s.", user.FirstName, "created", dMap[105].Action, dMap[105].Amount, dMap[105].Name))
}

//
// Test get Activities 02 - limit
//
func TestGetActivities02(t *testing.T) {
	// Data map
	dMap := make(map[uint]models.Activity)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test user
	user := test.GetRandomUser(33)
	db.Save(&user)

	// Create like 105 ledger entries. This will create Activities
	for i := 0; i < 105; i++ {
		l := test.GetRandomLedger(33)
		db.LedgerCreate(&l)

		// Set the ledger type
		ledgerType := "expense"

		if l.Amount > 0 {
			ledgerType = "income"
		}

		// Get the contact name.
		contactName := l.Contact.Name

		if len(contactName) == 0 {
			contactName = l.Contact.FirstName + " " + l.Contact.LastName
		}

		// Add to the activity log
		y := models.Activity{
			AccountId: uint(33),
			UserId:    user.Id,
			Action:    ledgerType,
			SubAction: "create",
			Name:      contactName,
			Amount:    l.Amount,
			LedgerId:  l.Id,
		}

		db.Create(&y)

		dMap[l.Id] = y
	}

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/activities?limit=25", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.GET("/api/v3/:account/activities", c.GetActivities)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []models.Activity{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, len(results), 25)
	st.Expect(t, results[0].Id, uint(105))
	st.Expect(t, results[1].Id, uint(104))
	st.Expect(t, results[2].Id, uint(103))
	st.Expect(t, results[0].Message, fmt.Sprintf("%s %s an %s ledger entry of %.2f from %s.", user.FirstName, "created", dMap[105].Action, dMap[105].Amount, dMap[105].Name))
}

//
// Test get Activities 03 - limit / page
//
func TestGetActivities03(t *testing.T) {
	// Data map
	dMap := make(map[uint]models.Activity)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test user
	user := test.GetRandomUser(33)
	db.Save(&user)

	// Create like 105 ledger entries. This will create Activities
	for i := 0; i < 105; i++ {
		l := test.GetRandomLedger(33)
		db.LedgerCreate(&l)

		// Set the ledger type
		ledgerType := "expense"

		if l.Amount > 0 {
			ledgerType = "income"
		}

		// Get the contact name.
		contactName := l.Contact.Name

		if len(contactName) == 0 {
			contactName = l.Contact.FirstName + " " + l.Contact.LastName
		}

		// Add to the activity log
		y := models.Activity{
			AccountId: uint(33),
			UserId:    user.Id,
			Action:    ledgerType,
			SubAction: "create",
			Name:      contactName,
			Amount:    l.Amount,
			LedgerId:  l.Id,
		}

		db.Create(&y)

		dMap[l.Id] = y
	}

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/activities?limit=25&page=2", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.GET("/api/v3/:account/activities", c.GetActivities)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []models.Activity{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, len(results), 25)
	st.Expect(t, results[0].Id, uint(80))
	st.Expect(t, results[1].Id, uint(79))
	st.Expect(t, results[2].Id, uint(78))
	st.Expect(t, results[0].Message, fmt.Sprintf("%s %s an %s ledger entry of %.2f from %s.", user.FirstName, "created", dMap[80].Action, dMap[80].Amount, dMap[80].Name))
}

/* End File */
