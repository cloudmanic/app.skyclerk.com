//
// Date: 2019-01-13
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2019-01-13
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
	"github.com/tidwall/gjson"

	"github.com/cloudmanic/skyclerk.com/models"
)

//
// TestGetContacts01 Test get contacts 01
//
func TestGetContacts01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Test labels. -- First 2 are to make sure we don't get them as they are not our account.
	db.Save(&models.Contact{AccountId: 34, Name: "Apple Inc.", FirstName: "", LastName: ""})
	db.Save(&models.Contact{AccountId: 34, Name: "Matthews Etc. LLC", FirstName: "Spicer", LastName: "Matthews"})
	db.Save(&models.Contact{AccountId: 33, Name: "", FirstName: "Jane", LastName: "Wells"})
	db.Save(&models.Contact{AccountId: 33, Name: "", FirstName: "Mike", LastName: "Rosso"})
	db.Save(&models.Contact{AccountId: 33, Name: "Zoo Inc.", FirstName: "", LastName: ""})
	db.Save(&models.Contact{AccountId: 33, Name: "Abc Inc.", FirstName: "Katie", LastName: "Matthews"})
	db.Save(&models.Contact{AccountId: 33, Name: "Dope Dealer, LLC", FirstName: "", LastName: ""})

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/contacts", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/contacts", c.GetContacts)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []models.Contact{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, results[0].Id, uint(3))
	st.Expect(t, results[1].Id, uint(4))
	st.Expect(t, results[2].Id, uint(6))
	st.Expect(t, results[3].Id, uint(7))
	st.Expect(t, results[4].Id, uint(5))
	st.Expect(t, results[0].Name, "")
	st.Expect(t, results[1].Name, "")
	st.Expect(t, results[2].Name, "Abc Inc.")
	st.Expect(t, results[3].Name, "Dope Dealer, LLC")
	st.Expect(t, results[4].Name, "Zoo Inc.")
}

//
// TestGetContacts02 Test get contacts 02
//
func TestGetContacts02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Test labels. -- First 2 are to make sure we don't get them as they are not our account.
	db.Save(&models.Contact{AccountId: 34, Name: "Apple Inc.", FirstName: "", LastName: ""})
	db.Save(&models.Contact{AccountId: 34, Name: "Matthews Etc. LLC", FirstName: "Spicer", LastName: "Matthews"})
	db.Save(&models.Contact{AccountId: 33, Name: "", FirstName: "Jane", LastName: "Wells"})
	db.Save(&models.Contact{AccountId: 33, Name: "", FirstName: "Mike", LastName: "Rosso"})
	db.Save(&models.Contact{AccountId: 33, Name: "Zoo Inc.", FirstName: "", LastName: ""})
	db.Save(&models.Contact{AccountId: 33, Name: "Abc Inc.", FirstName: "Katie", LastName: "Matthews"})
	db.Save(&models.Contact{AccountId: 33, Name: "Dope Dealer, LLC", FirstName: "", LastName: ""})

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/contacts?sort=desc", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/contacts", c.GetContacts)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []models.Contact{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, results[0].Id, uint(5))
	st.Expect(t, results[1].Id, uint(7))
	st.Expect(t, results[2].Id, uint(6))
	st.Expect(t, results[3].Id, uint(3))
	st.Expect(t, results[4].Id, uint(4))
	st.Expect(t, results[0].Name, "Zoo Inc.")
	st.Expect(t, results[1].Name, "Dope Dealer, LLC")
	st.Expect(t, results[2].Name, "Abc Inc.")
	st.Expect(t, results[3].Name, "")
	st.Expect(t, results[4].Name, "")
}

//
// TestGetContact0 Test get contacts 01
//
func TestGetContact01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Test contacts. -- First 2 are to make sure we don't get them as they are not our account.
	db.Save(&models.Contact{AccountId: 34, Name: "Apple Inc.", FirstName: "", LastName: ""})
	db.Save(&models.Contact{AccountId: 34, Name: "Matthews Etc. LLC", FirstName: "Spicer", LastName: "Matthews"})
	db.Save(&models.Contact{AccountId: 33, Name: "", FirstName: "Jane", LastName: "Wells"})
	db.Save(&models.Contact{AccountId: 33, Name: "", FirstName: "Mike", LastName: "Rosso"})
	db.Save(&models.Contact{AccountId: 33, Name: "Zoo Inc.", FirstName: "", LastName: ""})
	db.Save(&models.Contact{AccountId: 33, Name: "Abc Inc.", FirstName: "Katie", LastName: "Matthews"})
	db.Save(&models.Contact{AccountId: 33, Name: "Dope Dealer, LLC", FirstName: "", LastName: ""})

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/contacts/4", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/contacts/:id", c.GetContact)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Contact{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, result.Id, uint(4))
	st.Expect(t, result.FirstName, "Mike")
	st.Expect(t, result.LastName, "Rosso")
}

//
// Test create Contact 01
//
func TestCreateContact01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Post data
	post := models.Contact{Name: "Contact #1"}

	// Get JSON
	postStr, _ := json.Marshal(post)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/contacts", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v3/33/contacts", c.CreateContact)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Contact{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 201)
	st.Expect(t, result.Name, "Contact #1")
	st.Expect(t, result.FirstName, "")
	st.Expect(t, result.LastName, "")

	// Double check the db.
	contact := models.Contact{}
	db.First(&contact, 1)
	st.Expect(t, contact.Id, uint(1))
	st.Expect(t, contact.Name, "Contact #1")
	st.Expect(t, contact.FirstName, "")
	st.Expect(t, contact.LastName, "")
}

//
// Test create Contact 02 - Duplicate Name
//
func TestCreateContact02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test contact to conflict with
	db.Save(&models.Contact{AccountId: 33, Name: "Contact #1"})

	// Post data
	post := models.Contact{Name: "Contact #1"}

	// Get JSON
	postStr, _ := json.Marshal(post)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/contacts", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v3/33/contacts", c.CreateContact)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "errors.name").String(), "Contact name is already in use.")
}

//
// Test create Contact 03 - Duplicate First / Last
//
func TestCreateContact03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test contact to conflict with
	db.Save(&models.Contact{AccountId: 33, Name: "", FirstName: "Jane", LastName: "Wells"})

	// Post data
	post := models.Contact{FirstName: "Jane", LastName: "Wells"}

	// Get JSON
	postStr, _ := json.Marshal(post)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/contacts", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v3/33/contacts", c.CreateContact)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "errors.name").String(), "Contact first and last name is already in use.")
}

//
// Test create Contact 04 - Add first / last (with spaces)
//
func TestCreateContact04(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Post data
	post := models.Contact{FirstName: " Jane ", LastName: "   Wells "}

	// Get JSON
	postStr, _ := json.Marshal(post)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/contacts", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v3/33/contacts", c.CreateContact)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Contact{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 201)
	st.Expect(t, result.Name, "")
	st.Expect(t, result.FirstName, "Jane")
	st.Expect(t, result.LastName, "Wells")

	// Double check the db.
	contact := models.Contact{}
	db.First(&contact, 1)
	st.Expect(t, contact.Id, uint(1))
	st.Expect(t, contact.Name, "")
	st.Expect(t, contact.FirstName, "Jane")
	st.Expect(t, contact.LastName, "Wells")
}

//
// Test create Contact 05 - Full object test
//
func TestCreateContact05(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Post data
	post := models.Contact{
		Name:          "ABC Inc.",
		FirstName:     "Jane",
		LastName:      "Wells",
		Address:       "123 West Main Street",
		City:          "Newberg",
		State:         "OR",
		Zip:           "13601",
		Phone:         "543-876-9872",
		Fax:           "777-888-666",
		Website:       "https://google.com",
		AccountNumber: "1234abc",
		Email:         "jane@gmail.com",
		Twitter:       "@spicermatthews",
		Facebook:      "http://www.facebook.com/user-profile",
		Linkedin:      "http://www.linkedin.com/your-profile",
		Country:       "United States",
	}

	// Get JSON
	postStr, _ := json.Marshal(post)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/contacts", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v3/33/contacts", c.CreateContact)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Contact{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 201)
	st.Expect(t, result.Name, "ABC Inc.")
	st.Expect(t, result.FirstName, "Jane")
	st.Expect(t, result.LastName, "Wells")
	st.Expect(t, result.Address, "123 West Main Street")
	st.Expect(t, result.City, "Newberg")
	st.Expect(t, result.State, "OR")
	st.Expect(t, result.Zip, "13601")
	st.Expect(t, result.Country, "United States")
	st.Expect(t, result.Phone, "543-876-9872")
	st.Expect(t, result.Fax, "777-888-666")
	st.Expect(t, result.Website, "https://google.com")
	st.Expect(t, result.AccountNumber, "1234abc")
	st.Expect(t, result.Email, "jane@gmail.com")
	st.Expect(t, result.Twitter, "@spicermatthews")
	st.Expect(t, result.Facebook, "http://www.facebook.com/user-profile")
	st.Expect(t, result.Linkedin, "http://www.linkedin.com/your-profile")

	// Double check the db.
	contact := models.Contact{}
	db.First(&contact, 1)
	st.Expect(t, contact.Id, uint(1))
	st.Expect(t, contact.Name, "ABC Inc.")
	st.Expect(t, contact.FirstName, "Jane")
	st.Expect(t, contact.LastName, "Wells")
	st.Expect(t, contact.Address, "123 West Main Street")
	st.Expect(t, contact.City, "Newberg")
	st.Expect(t, contact.State, "OR")
	st.Expect(t, contact.Zip, "13601")
	st.Expect(t, contact.Country, "United States")
	st.Expect(t, contact.Phone, "543-876-9872")
	st.Expect(t, contact.Fax, "777-888-666")
	st.Expect(t, contact.Website, "https://google.com")
	st.Expect(t, contact.AccountNumber, "1234abc")
	st.Expect(t, contact.Email, "jane@gmail.com")
	st.Expect(t, contact.Twitter, "@spicermatthews")
	st.Expect(t, contact.Facebook, "http://www.facebook.com/user-profile")
	st.Expect(t, contact.Linkedin, "http://www.linkedin.com/your-profile")
}

//
// TestUpdateContact01 - Test update contact 01
//
func TestUpdateContact01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Test contacts. -- First 2 are to make sure we don't get them as they are not our account.
	db.Save(&models.Contact{AccountId: 34, Name: "Apple Inc.", FirstName: "", LastName: ""})
	db.Save(&models.Contact{AccountId: 34, Name: "Matthews Etc. LLC", FirstName: "Spicer", LastName: "Matthews"})
	db.Save(&models.Contact{AccountId: 33, Name: "", FirstName: "Jane", LastName: "Wells"})
	db.Save(&models.Contact{AccountId: 33, Name: "", FirstName: "Mike", LastName: "Rosso"})
	db.Save(&models.Contact{AccountId: 33, Name: "Zoo Inc.", FirstName: "", LastName: ""})
	db.Save(&models.Contact{AccountId: 33, Name: "Abc Inc.", FirstName: "Katie", LastName: "Matthews"})
	db.Save(&models.Contact{AccountId: 33, Name: "Dope Dealer, LLC", FirstName: "", LastName: ""})

	// Put data
	conPut := models.Contact{Name: "Wells Holdings, LLC", FirstName: "Jane", LastName: "Wells"}

	// Get JSON
	putStr, _ := json.Marshal(conPut)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/contacts/3", bytes.NewBuffer(putStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.PUT("/api/v3/:account/contacts/:id", c.UpdateContact)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Contact{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, result.Name, "Wells Holdings, LLC")
	st.Expect(t, result.FirstName, "Jane")
	st.Expect(t, result.LastName, "Wells")

	// Double check the db.
	con := models.Contact{}
	db.First(&con, 3)
	st.Expect(t, con.Id, uint(3))
	st.Expect(t, con.Name, "Wells Holdings, LLC")
	st.Expect(t, con.FirstName, "Jane")
	st.Expect(t, con.LastName, "Wells")
}

//
// TestUpdateContact02 - Test update contact 02
//
func TestUpdateContact02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Test contacts. -- First 2 are to make sure we don't get them as they are not our account.
	db.Save(&models.Contact{AccountId: 34, Name: "Apple Inc.", FirstName: "", LastName: ""})
	db.Save(&models.Contact{AccountId: 34, Name: "Matthews Etc. LLC", FirstName: "Spicer", LastName: "Matthews"})
	db.Save(&models.Contact{AccountId: 33, Name: "", FirstName: "Jane", LastName: "Wells"})
	db.Save(&models.Contact{AccountId: 33, Name: "", FirstName: "Mike", LastName: "Rosso"})
	db.Save(&models.Contact{AccountId: 33, Name: "Zoo Inc.", FirstName: "", LastName: ""})
	db.Save(&models.Contact{AccountId: 33, Name: "Abc Inc.", FirstName: "Katie", LastName: "Matthews"})
	db.Save(&models.Contact{AccountId: 33, Name: "Dope Dealer, LLC", FirstName: "", LastName: ""})

	// Put data
	conPut := models.Contact{Name: "Abc Inc.", FirstName: "Katie", LastName: "Matthews"}

	// Get JSON
	putStr, _ := json.Marshal(conPut)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/contacts/3", bytes.NewBuffer(putStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.PUT("/api/v3/:account/contacts/:id", c.UpdateContact)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"errors":{"name":"Contact company name, first, and last name is already in use."}}`)
}

//
// TestDeleteContact01 - Test delete Contact 01
//
func TestDeleteContact01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Test contacts. -- First 2 are to make sure we don't get them as they are not our account.
	db.Save(&models.Contact{AccountId: 34, Name: "Apple Inc.", FirstName: "", LastName: ""})
	db.Save(&models.Contact{AccountId: 34, Name: "Matthews Etc. LLC", FirstName: "Spicer", LastName: "Matthews"})
	db.Save(&models.Contact{AccountId: 33, Name: "", FirstName: "Jane", LastName: "Wells"})
	db.Save(&models.Contact{AccountId: 33, Name: "", FirstName: "Mike", LastName: "Rosso"})
	db.Save(&models.Contact{AccountId: 33, Name: "Zoo Inc.", FirstName: "", LastName: ""})
	db.Save(&models.Contact{AccountId: 33, Name: "Abc Inc.", FirstName: "Katie", LastName: "Matthews"})
	db.Save(&models.Contact{AccountId: 33, Name: "Dope Dealer, LLC", FirstName: "", LastName: ""})

	// Make sure our delete record is in the DB
	con1 := models.Contact{}
	db.Find(&con1, 5)
	st.Expect(t, con1.Id, uint(5))

	// Setup request
	req, _ := http.NewRequest("DELETE", "/api/v3/33/contacts/5", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.DELETE("/api/v3/:account/contacts/:id", c.DeleteContact)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 204)

	// Double check the db.
	cons := []models.Contact{}
	db.Find(&cons)
	st.Expect(t, len(cons), 6)
	st.Expect(t, cons[3].Id, uint(4))
	st.Expect(t, cons[3].Name, "")
	st.Expect(t, cons[4].Id, uint(6))
	st.Expect(t, cons[4].Name, "Abc Inc.")
	st.Expect(t, cons[5].Id, uint(7))
	st.Expect(t, cons[5].Name, "Dope Dealer, LLC")

	// More double check
	con := models.Contact{}
	db.Find(&con, 5)
	st.Expect(t, con.Id, uint(0))
}

//
// TestDeleteContact02 - Already part of a ledger.
//
func TestDeleteContact02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Test contacts. -- First 2 are to make sure we don't get them as they are not our account.
	db.Save(&models.Contact{AccountId: 34, Name: "Apple Inc.", FirstName: "", LastName: ""})
	db.Save(&models.Contact{AccountId: 34, Name: "Matthews Etc. LLC", FirstName: "Spicer", LastName: "Matthews"})
	db.Save(&models.Contact{AccountId: 33, Name: "", FirstName: "Jane", LastName: "Wells"})
	db.Save(&models.Contact{AccountId: 33, Name: "", FirstName: "Mike", LastName: "Rosso"})
	db.Save(&models.Contact{AccountId: 33, Name: "Zoo Inc.", FirstName: "", LastName: ""})
	db.Save(&models.Contact{AccountId: 33, Name: "Abc Inc.", FirstName: "Katie", LastName: "Matthews"})
	db.Save(&models.Contact{AccountId: 33, Name: "Dope Dealer, LLC", FirstName: "", LastName: ""})

	// Create test ledger
	db.Save(&models.Ledger{AccountId: 33, ContactId: 5})

	// Setup request
	req, _ := http.NewRequest("DELETE", "/api/v3/33/contacts/5", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.DELETE("/api/v3/:account/contacts/:id", c.DeleteContact)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "error").String(), "Can not delete contact. It is in use by a ledger entry.")
}

//
// TestDeleteContact03 - Failed account
//
func TestDeleteContact03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Test contacts. -- First 2 are to make sure we don't get them as they are not our account.
	db.Save(&models.Contact{AccountId: 34, Name: "Apple Inc.", FirstName: "", LastName: ""})
	db.Save(&models.Contact{AccountId: 34, Name: "Matthews Etc. LLC", FirstName: "Spicer", LastName: "Matthews"})
	db.Save(&models.Contact{AccountId: 33, Name: "", FirstName: "Jane", LastName: "Wells"})
	db.Save(&models.Contact{AccountId: 33, Name: "", FirstName: "Mike", LastName: "Rosso"})
	db.Save(&models.Contact{AccountId: 33, Name: "Zoo Inc.", FirstName: "", LastName: ""})
	db.Save(&models.Contact{AccountId: 33, Name: "Abc Inc.", FirstName: "Katie", LastName: "Matthews"})
	db.Save(&models.Contact{AccountId: 33, Name: "Dope Dealer, LLC", FirstName: "", LastName: ""})

	// Setup request
	req, _ := http.NewRequest("DELETE", "/api/v3/33/contacts/1", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.DELETE("/api/v3/:account/contacts/:id", c.DeleteContact)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "error").String(), "Contact not found.")
}

/* End File */
