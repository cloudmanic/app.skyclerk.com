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

	"github.com/cloudmanic/skyclerk.com/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
	"github.com/tidwall/gjson"
)

//
// Test create Contact 01
//
func TestCreateContact01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Post data
	post := models.Contact{Name: "Contact #1"}

	// Get JSON
	postStr, _ := json.Marshal(post)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v1/33/contacts", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("account", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v1/33/contacts", c.CreateContact)
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
	db, _ := models.NewDB()
	defer db.Close()

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
	req, _ := http.NewRequest("POST", "/api/v1/33/contacts", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("account", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v1/33/contacts", c.CreateContact)
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
	db, _ := models.NewDB()
	defer db.Close()

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
	req, _ := http.NewRequest("POST", "/api/v1/33/contacts", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("account", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v1/33/contacts", c.CreateContact)
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
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Post data
	post := models.Contact{FirstName: " Jane ", LastName: "   Wells "}

	// Get JSON
	postStr, _ := json.Marshal(post)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v1/33/contacts", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("account", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v1/33/contacts", c.CreateContact)
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
	db, _ := models.NewDB()
	defer db.Close()

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
	req, _ := http.NewRequest("POST", "/api/v1/33/contacts", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("account", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v1/33/contacts", c.CreateContact)
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

/* End File */