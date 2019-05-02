//
// Date: 2019-03-14
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
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
	"time"

	"app.skyclerk.com/backend/library/test"
	"app.skyclerk.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// TestGetLedgers01 Test get ledgers 01
//
func TestGetLedgers01(t *testing.T) {
	// Data map
	dMap := make(map[uint]models.Ledger)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create like 10 ledger entries. Diffent account.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(23)
		db.LedgerCreate(&l)
	}

	// Create like 105 ledger entries.
	for i := 0; i < 105; i++ {
		l := test.GetRandomLedger(33)
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Create like 10 ledger entries. Diffent account.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(43)
		db.LedgerCreate(&l)
	}

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/ledger", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/ledger", c.GetLedgers)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []models.Ledger{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, len(results), 50)
	st.Expect(t, w.HeaderMap["X-Offset"][0], "0")
	st.Expect(t, w.HeaderMap["X-Limit"][0], "50")
	st.Expect(t, w.HeaderMap["X-No-Limit-Count"][0], "105")
	st.Expect(t, w.HeaderMap["X-Last-Page"][0], "false")

	for key, row := range results {
		st.Expect(t, row.Id, dMap[row.Id].Id)
		st.Expect(t, row.AccountId, uint(33))
		st.Expect(t, row.Date.Format("2006-01-02"), dMap[row.Id].Date.Format("2006-01-02"))
		st.Expect(t, row.Amount, dMap[row.Id].Amount)
		st.Expect(t, row.Note, dMap[row.Id].Note)
		st.Expect(t, row.Contact.Name, dMap[row.Id].Contact.Name)
		st.Expect(t, row.Contact.FirstName, dMap[row.Id].Contact.FirstName)
		st.Expect(t, row.Contact.LastName, dMap[row.Id].Contact.LastName)
		st.Expect(t, row.Contact.Email, dMap[row.Id].Contact.Email)
		st.Expect(t, row.Contact.AccountId, uint(33))
		st.Expect(t, row.Category.AccountId, uint(33))
		st.Expect(t, row.Category.Name, dMap[row.Id].Category.Name)
		st.Expect(t, row.Labels[0].AccountId, uint(33))

		// Verfiy default Order
		if key > 0 {
			diff := row.Date.Sub(results[key-1].Date)
			st.Expect(t, (diff <= 0), true)
		}
	}

	// ----------- Test Paging 2 ---------- //

	// Setup request
	req2, _ := http.NewRequest("GET", "/api/v3/33/ledger?page=2", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	// Grab result and convert to strut
	results2 := []models.Ledger{}
	err = json.Unmarshal([]byte(w2.Body.String()), &results2)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w2.Code, 200)
	st.Expect(t, len(results2), 50)
	st.Expect(t, w2.HeaderMap["X-Offset"][0], "50")
	st.Expect(t, w2.HeaderMap["X-Limit"][0], "50")
	st.Expect(t, w2.HeaderMap["X-No-Limit-Count"][0], "105")
	st.Expect(t, w2.HeaderMap["X-Last-Page"][0], "false")

	// ----------- Test Paging 3 ---------- //

	// Setup request
	req3, _ := http.NewRequest("GET", "/api/v3/33/ledger?page=3", nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	// Grab result and convert to strut
	results3 := []models.Ledger{}
	err = json.Unmarshal([]byte(w3.Body.String()), &results3)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w3.Code, 200)
	st.Expect(t, len(results3), 5)
	st.Expect(t, w3.HeaderMap["X-Offset"][0], "100")
	st.Expect(t, w3.HeaderMap["X-Limit"][0], "50")
	st.Expect(t, w3.HeaderMap["X-No-Limit-Count"][0], "105")
	st.Expect(t, w3.HeaderMap["X-Last-Page"][0], "true")
}

//
// TestGetLedgers02 Test type filters
//
func TestGetLedgers02(t *testing.T) {
	// Data map
	dMap := make(map[uint]models.Ledger)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create like 10 ledger entries. Diffent account.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(23)
		db.LedgerCreate(&l)
	}

	// Create like 105 ledger entries.
	for i := 0; i < 105; i++ {
		l := test.GetRandomLedger(33)
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Create like 10 ledger entries. Diffent account.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(43)
		db.LedgerCreate(&l)
	}

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/ledger?type=income", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/ledger", c.GetLedgers)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []models.Ledger{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, w.HeaderMap["X-Offset"][0], "0")
	st.Expect(t, w.HeaderMap["X-Limit"][0], "50")
	st.Expect(t, w.HeaderMap["X-Last-Page"][0], "false")

	for key, row := range results {
		st.Expect(t, (row.Amount > 0), true)

		// Verfiy default Order
		if key > 0 {
			diff := row.Date.Sub(results[key-1].Date)
			st.Expect(t, (diff <= 0), true)
		}
	}

	// ----------- Test expense ---------- //

	// Setup request
	req2, _ := http.NewRequest("GET", "/api/v3/33/ledger?type=expense", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	// Grab result and convert to strut
	results2 := []models.Ledger{}
	err = json.Unmarshal([]byte(w2.Body.String()), &results2)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, w.HeaderMap["X-Offset"][0], "0")
	st.Expect(t, w.HeaderMap["X-Limit"][0], "50")
	st.Expect(t, w.HeaderMap["X-Last-Page"][0], "false")

	for key, row := range results2 {
		st.Expect(t, (row.Amount < 0), true)

		// Verfiy default Order
		if key > 0 {
			diff := row.Date.Sub(results2[key-1].Date)
			st.Expect(t, (diff <= 0), true)
		}
	}
}

//
// TestGetLedgers03 Test type filters
//
func TestGetLedgers03(t *testing.T) {
	// Data map
	dMap := make(map[uint]models.Ledger)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create like 10 ledger entries. Diffent account.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(23)
		db.LedgerCreate(&l)
	}

	// Create like 105 ledger entries.
	for i := 0; i < 105; i++ {
		l := test.GetRandomLedger(33)
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Create like 10 ledger entries. Diffent account.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(43)
		db.LedgerCreate(&l)
	}

	// Setup request
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v3/33/ledger?category_id=%d", dMap[106].CategoryId), nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/ledger", c.GetLedgers)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []models.Ledger{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, w.HeaderMap["X-Offset"][0], "0")
	st.Expect(t, w.HeaderMap["X-Limit"][0], "50")

	for key, row := range results {
		st.Expect(t, (dMap[106].CategoryId == row.Category.Id), true)

		// Verfiy default Order
		if key > 0 {
			diff := row.Date.Sub(results[key-1].Date)
			st.Expect(t, (diff <= 0), true)
		}
	}
}

//
// TestGetLedgers04 Test label filters
//
func TestGetLedgers04(t *testing.T) {
	// Data map
	dMap := make(map[uint]models.Ledger)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create like 10 ledger entries. Diffent account.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(23)
		db.LedgerCreate(&l)
	}

	// Create like 105 ledger entries.
	for i := 0; i < 105; i++ {
		l := test.GetRandomLedger(33)
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Create like 10 ledger entries. Diffent account.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(43)
		db.LedgerCreate(&l)
	}

	// Setup request
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v3/33/ledger?label_ids=%d,%d", dMap[106].Labels[0].Id, dMap[106].Labels[1].Id), nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/ledger", c.GetLedgers)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []models.Ledger{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, w.HeaderMap["X-Offset"][0], "0")
	st.Expect(t, w.HeaderMap["X-Limit"][0], "50")

	for key, row := range results {
		found := false
		for _, row2 := range row.Labels {
			if (row2.Id == dMap[106].Labels[0].Id) || (row2.Id == dMap[106].Labels[1].Id) {
				found = true
			}
		}
		st.Expect(t, found, true)

		// Verfiy default Order
		if key > 0 {
			diff := row.Date.Sub(results[key-1].Date)
			st.Expect(t, (diff <= 0), true)
		}
	}
}

//
// TestGetLedgers05 Test year filters
//
func TestGetLedgers05(t *testing.T) {
	// Data map
	dMap := make(map[uint]models.Ledger)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("testing_db")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create like 10 ledger entries. Diffent account.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(23)
		db.LedgerCreate(&l)
	}

	// Create like 105 ledger entries.
	for i := 0; i < 105; i++ {
		l := test.GetRandomLedger(33)
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Create like 10 ledger entries. Diffent account.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(43)
		db.LedgerCreate(&l)
	}

	// Setup request
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v3/33/ledger?year=2017"), nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/ledger", c.GetLedgers)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []models.Ledger{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, w.HeaderMap["X-Offset"][0], "0")
	st.Expect(t, w.HeaderMap["X-Limit"][0], "50")

	for key, row := range results {
		st.Expect(t, row.Date.Format("2006"), "2017")

		// Verfiy default Order
		if key > 0 {
			diff := row.Date.Sub(results[key-1].Date)
			st.Expect(t, (diff <= 0), true)
		}
	}
}

//
// TestGetLedger01 Test get ledger 01
//
func TestGetLedger01(t *testing.T) {
	// Data map
	dMap := make(map[uint]models.Ledger)

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create like 10 ledger entries. Diffent account.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(23)
		db.LedgerCreate(&l)
	}

	// Create like 10 ledger entries.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		db.LedgerCreate(&l)
		dMap[l.Id] = l
	}

	// Create like 10 ledger entries. Diffent account.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(43)
		db.LedgerCreate(&l)
	}

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/ledger/15", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/ledger/:id", c.GetLedger)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Ledger{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, result.Id, dMap[result.Id].Id)
	st.Expect(t, result.AccountId, uint(33))
	st.Expect(t, result.Date.Format("2006-01-02"), dMap[result.Id].Date.Format("2006-01-02"))
	st.Expect(t, result.Amount, dMap[result.Id].Amount)
	st.Expect(t, result.Note, dMap[result.Id].Note)
	st.Expect(t, result.Contact.Name, dMap[result.Id].Contact.Name)
	st.Expect(t, result.Contact.FirstName, dMap[result.Id].Contact.FirstName)
	st.Expect(t, result.Contact.LastName, dMap[result.Id].Contact.LastName)
	st.Expect(t, result.Contact.Email, dMap[result.Id].Contact.Email)
	st.Expect(t, result.Contact.AccountId, uint(33))
	st.Expect(t, result.Category.AccountId, uint(33))
	st.Expect(t, result.Category.Name, dMap[result.Id].Category.Name)
	st.Expect(t, result.Labels[0].AccountId, uint(33))
}

//
// TestGetLedger02 Test get ledger 02 - No perms
//
func TestGetLedger02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create like 10 ledger entries. Diffent account.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(23)
		db.LedgerCreate(&l)
	}

	// Create like 10 ledger entries.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		db.LedgerCreate(&l)
	}

	// Create like 10 ledger entries. Diffent account.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(43)
		db.LedgerCreate(&l)
	}

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/ledger/5", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/ledger/:id", c.GetLedger)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"Ledger entry not found."}`)
}

// -------------- Create Ledger ---------------------- //

//
// Test create Ledger 01
//
func TestCreateLedger01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Post data
	post := test.GetRandomLedger(33)

	// Get JSON
	postStr, _ := json.Marshal(post)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/ledger", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v3/33/ledger", c.CreateLedger)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Ledger{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 201)
	st.Expect(t, result.Id, uint(1))
	st.Expect(t, result.AccountId, uint(33))
	st.Expect(t, result.Date.Format("2006-01-02"), post.Date.Format("2006-01-02"))
	st.Expect(t, result.Amount, post.Amount)
	st.Expect(t, result.Note, post.Note)
	st.Expect(t, result.Contact.Name, post.Contact.Name)
	st.Expect(t, result.Contact.FirstName, post.Contact.FirstName)
	st.Expect(t, result.Contact.LastName, post.Contact.LastName)
	st.Expect(t, result.Contact.Email, post.Contact.Email)
	st.Expect(t, result.Contact.AccountId, uint(33))
	st.Expect(t, result.Category.AccountId, uint(33))
	st.Expect(t, result.Category.Name, post.Category.Name)
	st.Expect(t, result.Labels[0].AccountId, uint(33))

	// Double check the db.
	result1, err := db.GetLedgerByAccountAndId(uint(33), uint(1))
	st.Expect(t, err, nil)
	st.Expect(t, result1.Id, uint(1))
	st.Expect(t, result1.AccountId, uint(33))
	st.Expect(t, result1.Date.Format("2006-01-02"), post.Date.Format("2006-01-02"))
	st.Expect(t, result1.Amount, post.Amount)
	st.Expect(t, result1.Note, post.Note)
	st.Expect(t, result1.Contact.Name, post.Contact.Name)
	st.Expect(t, result1.Contact.FirstName, post.Contact.FirstName)
	st.Expect(t, result1.Contact.LastName, post.Contact.LastName)
	st.Expect(t, result1.Contact.Email, post.Contact.Email)
	st.Expect(t, result1.Contact.AccountId, uint(33))
	st.Expect(t, result1.Category.AccountId, uint(33))
	st.Expect(t, result1.Category.Name, post.Category.Name)
	st.Expect(t, result1.Labels[0].AccountId, uint(33))
}

//
// Test create Ledger 02 - Validation
//
func TestCreateLedger02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Get JSON
	postStr := []byte(`{ "amount": 88.12 }`)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/ledger", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v3/33/ledger", c.CreateLedger)
	r.ServeHTTP(w, req)

	// Check result
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"errors":{"category":"Category name is required.","contact":"Contact name is required.","date":"The date field is required."}}`)

	// ------------ Test 2 ---------------- //

	// Get JSON
	postStr = []byte(`{ "amount": 88.12, "date": "2018-08-02" }`)

	// Setup request
	req, _ = http.NewRequest("POST", "/api/v3/33/ledger", bytes.NewBuffer(postStr))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check result
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"Invalid JSON in body. There is a chance the JSON maybe valid but does not match the data type requirements. For example maybe you passed a string in for an integer."}`)

	// ------------ Test 3 ---------------- //

	// Get JSON
	postStr = []byte(`{ "amount": 88.12, "date": "2018-08-02T08:18:20Z", "category": { "name": "woots" } }`)

	// Setup request
	req, _ = http.NewRequest("POST", "/api/v3/33/ledger", bytes.NewBuffer(postStr))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check result
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"errors":{"category":"Category type is required.","contact":"Contact name is required."}}`)

	// ------------ Test 4 ---------------- //

	// Get JSON
	postStr = []byte(`{ "amount": 88.12, "date": "2018-08-02T08:18:20Z", "category": { "name": "woots" }, "contact": { "first_name": "Jane", "last_name": "Wells", "name": "ABC Inc." } }`)

	// Setup request
	req, _ = http.NewRequest("POST", "/api/v3/33/ledger", bytes.NewBuffer(postStr))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check result
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"errors":{"category":"Category type is required."}}`)
}

//
// Test create Ledger 03 - trim and more
//
func TestCreateLedger03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Get JSON
	postStr := []byte(`{ "amount": 88.12, "date": "2018-08-02T08:18:20Z", "category": { "name": "  woots  ", "type": "  2  " }, "contact": { "first_name": "   Jane  ", "last_name": "  Wells  ", "name": "  ABC Inc.  " } }`)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/ledger", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v3/33/ledger", c.CreateLedger)
	r.ServeHTTP(w, req)

	// Check result
	st.Expect(t, w.Code, 201)

	// Double check the db.
	result1, err := db.GetLedgerByAccountAndId(uint(33), uint(1))
	st.Expect(t, err, nil)
	st.Expect(t, result1.Id, uint(1))
	st.Expect(t, result1.AccountId, uint(33))
	st.Expect(t, result1.Date.Format("2006-01-02"), "2018-08-02")
	st.Expect(t, result1.Amount, 88.12)
	st.Expect(t, result1.Note, "")
	st.Expect(t, result1.Contact.Name, "ABC Inc.")
	st.Expect(t, result1.Contact.FirstName, "Jane")
	st.Expect(t, result1.Contact.LastName, "Wells")
	st.Expect(t, result1.Contact.AccountId, uint(33))
	st.Expect(t, result1.Category.AccountId, uint(33))
	st.Expect(t, result1.Category.Name, "woots")
}

//
// Test create Ledger 04 - exsiting Contact / Category
//
func TestCreateLedger04(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Save random contact
	con := test.GetRandomContact(33)
	db.Save(&con)

	// Save random category
	cat := test.GetRandomCategory(33)
	db.Save(&cat)

	// Get JSON
	postStr := []byte(`{ "amount": 88.12, "date": "2018-08-02T08:18:20Z", "category": { "name": "` + cat.Name + `", "type": "` + cat.Type + `" }, "contact": { "first_name": "` + con.FirstName + `", "last_name": "` + con.LastName + `", "name": "` + con.Name + `", "email": "` + con.Email + `" } }`)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/ledger", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v3/33/ledger", c.CreateLedger)
	r.ServeHTTP(w, req)

	// Check result
	st.Expect(t, w.Code, 201)

	// Double check the db.
	result1, err := db.GetLedgerByAccountAndId(uint(33), uint(1))
	st.Expect(t, err, nil)
	st.Expect(t, result1.Id, uint(1))
	st.Expect(t, result1.AccountId, uint(33))
	st.Expect(t, result1.Date.Format("2006-01-02"), "2018-08-02")
	st.Expect(t, result1.Amount, 88.12)
	st.Expect(t, result1.Note, "")
	st.Expect(t, result1.Contact.Id, uint(1))
	st.Expect(t, result1.Contact.Name, con.Name)
	st.Expect(t, result1.Contact.FirstName, con.FirstName)
	st.Expect(t, result1.Contact.LastName, con.LastName)
	st.Expect(t, result1.Contact.AccountId, uint(33))
	st.Expect(t, result1.Category.AccountId, uint(33))
	st.Expect(t, result1.Category.Id, uint(1))
	st.Expect(t, result1.Category.Name, cat.Name)
}

// -------------- Update Ledger ---------------------- //

//
// Test update Ledger 01
//
func TestUpdateLedger01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create like random ledger entries.
	l := test.GetRandomLedger(33)
	db.LedgerCreate(&l)

	// Udpate a few fields.
	l.Amount = 199.23
	l.Date = time.Date(2019, 01, 01, 17, 20, 01, 507451, time.UTC)
	l.Contact.FirstName = "Testfirst"
	l.Contact.LastName = "Testlast"
	l.Contact.Name = "Test Company"
	l.Category.Name = "Test Category"
	l.Labels = []models.Label{
		{Name: "Test Label Unit Test"},
	}

	// Get JSON
	postStr, _ := json.Marshal(l)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/ledger/1", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.PUT("/api/v3/33/ledger/:id", c.UpdateLedger)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Ledger{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, result.Id, uint(1))
	st.Expect(t, result.AccountId, uint(33))
	st.Expect(t, result.Date.Format("2006-01-02"), l.Date.Format("2006-01-02"))
	st.Expect(t, result.Amount, l.Amount)
	st.Expect(t, result.Note, l.Note)
	st.Expect(t, result.Contact.Name, l.Contact.Name)
	st.Expect(t, result.Contact.FirstName, l.Contact.FirstName)
	st.Expect(t, result.Contact.LastName, l.Contact.LastName)
	st.Expect(t, result.Contact.Email, l.Contact.Email)
	st.Expect(t, result.Contact.AccountId, uint(33))
	st.Expect(t, result.Category.AccountId, uint(33))
	st.Expect(t, result.Category.Name, l.Category.Name)
	st.Expect(t, result.Labels[0].AccountId, uint(33))
	st.Expect(t, result.Labels[0].Name, "Test Label Unit Test")
	st.Expect(t, len(result.Labels), 1)

	// Double check the db.
	result1, err := db.GetLedgerByAccountAndId(uint(33), uint(1))
	st.Expect(t, err, nil)
	st.Expect(t, result1.Id, uint(1))
	st.Expect(t, result1.AccountId, uint(33))
	st.Expect(t, result1.Date.Format("2006-01-02"), l.Date.Format("2006-01-02"))
	st.Expect(t, result1.Amount, l.Amount)
	st.Expect(t, result1.Note, l.Note)
	st.Expect(t, result1.Contact.Name, l.Contact.Name)
	st.Expect(t, result1.Contact.FirstName, l.Contact.FirstName)
	st.Expect(t, result1.Contact.LastName, l.Contact.LastName)
	st.Expect(t, result1.Contact.Email, l.Contact.Email)
	st.Expect(t, result1.Contact.AccountId, uint(33))
	st.Expect(t, result1.Category.AccountId, uint(33))
	st.Expect(t, result1.Category.Name, l.Category.Name)
	st.Expect(t, result1.Labels[0].AccountId, uint(33))
	st.Expect(t, result1.Labels[0].Name, "Test Label Unit Test")
	st.Expect(t, len(result1.Labels), 1)
}

//
// Test update Ledger 02 - No Perms
//
func TestUpdateLedger02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create like random ledger entries.
	l := test.GetRandomLedger(34)
	db.LedgerCreate(&l)

	// Udpate a few fields.
	l.Amount = 199.23
	l.Date = time.Date(2019, 01, 01, 17, 20, 01, 507451, time.UTC)
	l.Contact.FirstName = "Testfirst"
	l.Contact.LastName = "Testlast"
	l.Contact.Name = "Test Company"
	l.Category.Name = "Test Category"
	l.Labels = []models.Label{
		{Name: "Test Label Unit Test"},
	}

	// Get JSON
	postStr, _ := json.Marshal(l)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/ledger/1", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.PUT("/api/v3/33/ledger/:id", c.UpdateLedger)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"Ledger entry not found."}`)
}

//
// Test update Ledger 03 - Validation
//
func TestUpdateLedger03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create like random ledger entries.
	l := test.GetRandomLedger(33)
	db.LedgerCreate(&l)

	// Udpate a few fields.
	l.Amount = 199.23
	l.Date = time.Date(2019, 01, 01, 17, 20, 01, 507451, time.UTC)
	l.Contact.FirstName = ""
	l.Contact.LastName = ""
	l.Contact.Name = ""
	l.Category.Name = ""
	l.Labels = []models.Label{
		{Name: "Test Label Unit Test"},
	}

	// Get JSON
	postStr, _ := json.Marshal(l)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/ledger/1", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.PUT("/api/v3/33/ledger/:id", c.UpdateLedger)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"errors":{"category":"Category name is required.","contact":"Contact name is required."}}`)
}

// -------------- Delete Ledger ---------------------- //

//
// TestDeleteLedger01 - Test delete Ledger 01
//
func TestDeleteLedger01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create like 10 ledger entries.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		db.LedgerCreate(&l)
	}

	// Make sure our delete record is in the DB
	l1 := models.Ledger{}
	db.Find(&l1, 5)
	st.Expect(t, l1.Id, uint(5))

	// Setup request
	req, _ := http.NewRequest("DELETE", "/api/v3/33/ledger/5", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.DELETE("/api/v3/:account/ledger/:id", c.DeleteLedger)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 204)

	// More double check
	l := models.Ledger{}
	db.Find(&l, 5)
	st.Expect(t, l.Id, uint(0))

	// More double double check
	l2 := []models.Ledger{}
	db.Find(&l2)
	st.Expect(t, len(l2), 9)
}

//
// TestDeleteLedger02 - No perms
//
func TestDeleteLedger02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create like 10 ledger entries. Diffent account.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(23)
		db.LedgerCreate(&l)
	}

	// Create like 10 ledger entries.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(33)
		db.LedgerCreate(&l)
	}

	// Create like 10 ledger entries. Diffent account.
	for i := 0; i < 10; i++ {
		l := test.GetRandomLedger(43)
		db.LedgerCreate(&l)
	}

	// Setup request
	req, _ := http.NewRequest("DELETE", "/api/v3/33/ledger/5", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.DELETE("/api/v3/:account/ledger/:id", c.DeleteLedger)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"Ledger entry not found."}`)
}

/* End File */
