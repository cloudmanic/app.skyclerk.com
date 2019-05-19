//
// Date: 2018-03-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-29
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"app.skyclerk.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
	"github.com/tidwall/gjson"
)

//
// Test get a ledger 01
//
func TestGetLabels01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Test labels. -- First 2 are to make sure we don't get them as they are not our account.
	db.Save(&models.Label{AccountId: 34, Name: "No #1"})
	db.Save(&models.Label{AccountId: 34, Name: "No #2"})
	db.Save(&models.Label{AccountId: 33, Name: "label #1"})
	db.Save(&models.Label{AccountId: 33, Name: "label #2"})
	db.Save(&models.Label{AccountId: 33, Name: "label #3"})
	db.Save(&models.Label{AccountId: 33, Name: "Abc"})
	db.Save(&models.Label{AccountId: 33, Name: "Xyz"})

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/labels", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/labels", c.GetLabels)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []models.Label{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)

	st.Expect(t, results[0].Id, uint(6))
	st.Expect(t, results[1].Id, uint(3))
	st.Expect(t, results[2].Id, uint(4))
	st.Expect(t, results[3].Id, uint(5))
	st.Expect(t, results[4].Id, uint(7))

	st.Expect(t, results[0].Name, "Abc")
	st.Expect(t, results[1].Name, "label #1")
	st.Expect(t, results[2].Name, "label #2")
	st.Expect(t, results[3].Name, "label #3")
	st.Expect(t, results[4].Name, "Xyz")
}

//
// Test get a ledger 02 - Sort
//
func TestGetLabels02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Test labels.
	db.Save(&models.Label{AccountId: 33, Name: "label #1"})
	db.Save(&models.Label{AccountId: 33, Name: "label #2"})
	db.Save(&models.Label{AccountId: 33, Name: "label #3"})
	db.Save(&models.Label{AccountId: 33, Name: "Abc"})
	db.Save(&models.Label{AccountId: 33, Name: "Xyz"})

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/labels?sort=desc", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/labels", c.GetLabels)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []models.Label{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)

	st.Expect(t, results[0].Id, uint(5))
	st.Expect(t, results[1].Id, uint(3))
	st.Expect(t, results[2].Id, uint(2))
	st.Expect(t, results[3].Id, uint(1))
	st.Expect(t, results[4].Id, uint(4))

	st.Expect(t, results[0].Name, "Xyz")
	st.Expect(t, results[1].Name, "label #3")
	st.Expect(t, results[2].Name, "label #2")
	st.Expect(t, results[3].Name, "label #1")
	st.Expect(t, results[4].Name, "Abc")
}

//
// Test get a ledger 03 - Order
//
func TestGetLabels03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Test labels.
	db.Save(&models.Label{AccountId: 33, Name: "label #1"})
	db.Save(&models.Label{AccountId: 33, Name: "label #2"})
	db.Save(&models.Label{AccountId: 33, Name: "label #3"})
	db.Save(&models.Label{AccountId: 33, Name: "Abc"})
	db.Save(&models.Label{AccountId: 33, Name: "Xyz"})

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/labels?sort=desc&order=LabelsId", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/labels", c.GetLabels)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []models.Label{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)

	st.Expect(t, results[0].Id, uint(5))
	st.Expect(t, results[1].Id, uint(4))
	st.Expect(t, results[2].Id, uint(3))
	st.Expect(t, results[3].Id, uint(2))
	st.Expect(t, results[4].Id, uint(1))

	st.Expect(t, results[0].Name, "Xyz")
	st.Expect(t, results[1].Name, "Abc")
	st.Expect(t, results[2].Name, "label #3")
	st.Expect(t, results[3].Name, "label #2")
	st.Expect(t, results[4].Name, "label #1")
}

//
// Test get a ledger 04 - Order & Sort
//
func TestGetLabels04(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Test labels.
	db.Save(&models.Label{AccountId: 33, Name: "label #1"})
	db.Save(&models.Label{AccountId: 33, Name: "label #2"})
	db.Save(&models.Label{AccountId: 33, Name: "label #3"})
	db.Save(&models.Label{AccountId: 33, Name: "Abc"})
	db.Save(&models.Label{AccountId: 33, Name: "Xyz"})

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/labels?sort=asc&order=LabelsId", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/labels", c.GetLabels)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []models.Label{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)

	st.Expect(t, results[0].Id, uint(1))
	st.Expect(t, results[1].Id, uint(2))
	st.Expect(t, results[2].Id, uint(3))
	st.Expect(t, results[3].Id, uint(4))
	st.Expect(t, results[4].Id, uint(5))

	st.Expect(t, results[0].Name, "label #1")
	st.Expect(t, results[1].Name, "label #2")
	st.Expect(t, results[2].Name, "label #3")
	st.Expect(t, results[3].Name, "Abc")
	st.Expect(t, results[4].Name, "Xyz")
}

//
// Test get a ledger 05 - Failed Order Col
//
func TestGetLabels05(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Test labels.
	db.Save(&models.Label{AccountId: 33, Name: "label #1"})
	db.Save(&models.Label{AccountId: 33, Name: "label #2"})
	db.Save(&models.Label{AccountId: 33, Name: "label #3"})
	db.Save(&models.Label{AccountId: 33, Name: "Abc"})
	db.Save(&models.Label{AccountId: 33, Name: "Xyz"})

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/labels?sort=asc&order=FailedId", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/labels", c.GetLabels)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Body.String(), `{"error":"There was an error. Please contact help@skyclerk.com for help."}`)
}

//
// Test get a label 01
//
func TestGetLabel01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Test labels. -- First 2 are to make sure we don't get them as they are not our account.
	db.Save(&models.Label{AccountId: 34, Name: "No #1"})
	db.Save(&models.Label{AccountId: 34, Name: "No #2"})
	db.Save(&models.Label{AccountId: 33, Name: "label #1"})
	db.Save(&models.Label{AccountId: 33, Name: "label #2"})
	db.Save(&models.Label{AccountId: 33, Name: "label #3"})
	db.Save(&models.Label{AccountId: 33, Name: "Abc"})
	db.Save(&models.Label{AccountId: 33, Name: "Xyz"})

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/labels/3", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.GET("/api/v3/:account/labels/:id", c.GetLabel)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Label{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, result.Id, uint(3))
	st.Expect(t, result.Name, "label #1")
}

//
// Test get a label 02 - no perms
//
func TestGetLabel02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Test labels. -- First 2 are to make sure we don't get them as they are not our account.
	db.Save(&models.Label{AccountId: 34, Name: "No #1"})
	db.Save(&models.Label{AccountId: 34, Name: "No #2"})
	db.Save(&models.Label{AccountId: 33, Name: "label #1"})
	db.Save(&models.Label{AccountId: 33, Name: "label #2"})
	db.Save(&models.Label{AccountId: 33, Name: "label #3"})
	db.Save(&models.Label{AccountId: 33, Name: "Abc"})
	db.Save(&models.Label{AccountId: 33, Name: "Xyz"})

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/labels/2", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.GET("/api/v3/:account/labels/:id", c.GetLabel)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "error").String(), "Label not found.")
}

//
// Test create Label 01
//
func TestCreateLabel01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Post data
	lPost := models.Label{Name: "Label #1"}

	// Get JSON
	postStr, _ := json.Marshal(lPost)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/labels", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v3/:account/labels", c.CreateLabel)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Label{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 201)
	st.Expect(t, result.Name, "Label #1")

	// Double check the db.
	lb := models.Label{}
	db.First(&lb, 1)
	st.Expect(t, lb.Id, uint(1))
	st.Expect(t, lb.Name, "Label #1")
}

//
// Test create Label 02 - duplicate
//
func TestCreateLabel02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Set conflict
	db.Save(&models.Label{AccountId: 33, Name: "label #1"})

	// Post data
	lPost := models.Label{Name: "Label #1"}

	// Get JSON
	postStr, _ := json.Marshal(lPost)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/labels", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v3/:account/labels", c.CreateLabel)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "errors.name").String(), "Label name is already in use.")
}

//
// Test create Label 03 -- spaces
//
func TestCreateLabel03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Post data
	lPost := models.Label{Name: "  Label #1  "}

	// Get JSON
	postStr, _ := json.Marshal(lPost)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/labels", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v3/:account/labels", c.CreateLabel)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Label{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 201)
	st.Expect(t, result.Name, "Label #1")

	// Double check the db.
	lb := models.Label{}
	db.First(&lb, 1)
	st.Expect(t, lb.Id, uint(1))
	st.Expect(t, lb.Name, "Label #1")
}

//
// Test create Label 04 - duplicate / space
//
func TestCreateLabel04(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Set conflict
	db.Save(&models.Label{AccountId: 33, Name: "label #1"})

	// Post data
	lPost := models.Label{Name: "    Label #1    "}

	// Get JSON
	postStr, _ := json.Marshal(lPost)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/labels", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v3/:account/labels", c.CreateLabel)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "errors.name").String(), "Label name is already in use.")
}

//
// Test create Label 05 - blank
//
func TestCreateLabel05(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Post data
	lPost := models.Label{Name: ""}

	// Get JSON
	postStr, _ := json.Marshal(lPost)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/labels", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v3/:account/labels", c.CreateLabel)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "errors.name").String(), "The name field is required.")
}

//
// TestUpdateLabel01 - Test update Label 01
//
func TestUpdateLabel01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test label to conflict with
	db.Save(&models.Label{AccountId: 33, Name: "Label #1"})

	// Put data
	lbPut := models.Label{Name: "Label #1 Unit Test"}

	// Get JSON
	putStr, _ := json.Marshal(lbPut)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/labels/1", bytes.NewBuffer(putStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.PUT("/api/v3/:account/labels/:id", c.UpdateLabel)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Label{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, result.Name, "Label #1 Unit Test")

	// Double check the db.
	lb := models.Label{}
	db.First(&lb, 1)
	st.Expect(t, lb.Id, uint(1))
	st.Expect(t, lb.Name, "Label #1 Unit Test")
}

//
// TestUpdateCategory02 - Test update Label 02 - Duplicate Label
//
func TestUpdateLabel02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test label to conflict with
	db.Save(&models.Label{AccountId: 33, Name: "Label #1"})
	db.Save(&models.Label{AccountId: 33, Name: "Label #2"})

	// Put data
	lbPut := models.Label{Name: "Label #2"}

	// Get JSON
	putStr, _ := json.Marshal(lbPut)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/labels/1", bytes.NewBuffer(putStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.PUT("/api/v3/:account/labels/:id", c.UpdateLabel)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "errors.name").String(), "Label name is already in use.")
}

//
// TestUpdateLabel03 - Test update Label 03 - Duplicate same id
//
func TestUpdateLabel03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test label to conflict with
	db.Save(&models.Label{AccountId: 33, Name: "Label #1"})
	db.Save(&models.Label{AccountId: 33, Name: "Label #2"})

	// Put data
	lbPut := models.Label{Name: "Label #2"}

	// Get JSON
	putStr, _ := json.Marshal(lbPut)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/labels/2", bytes.NewBuffer(putStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.PUT("/api/v3/:account/labels/:id", c.UpdateLabel)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Label{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, result.Name, "Label #2")

	// Double check the db.
	lb := models.Label{}
	db.First(&lb, 2)
	st.Expect(t, lb.Id, uint(2))
	st.Expect(t, lb.Name, "Label #2")
}

//
// TestUpdateLabel04 - Test update Label 02 - Duplicate Label space
//
func TestUpdateLabel04(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test label to conflict with
	db.Save(&models.Label{AccountId: 33, Name: "Label #1"})
	db.Save(&models.Label{AccountId: 33, Name: "Label #2"})

	// Put data
	lbPut := models.Label{Name: "    Label #2    "}

	// Get JSON
	putStr, _ := json.Marshal(lbPut)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/labels/1", bytes.NewBuffer(putStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.PUT("/api/v3/:account/labels/:id", c.UpdateLabel)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "errors.name").String(), "Label name is already in use.")
}

//
// TestUpdateLabel05 - Test update Label 05 - Space
//
func TestUpdateLabel05(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test label to conflict with
	db.Save(&models.Label{AccountId: 33, Name: "Label #1"})

	// Put data
	lbPut := models.Label{Name: "     Label #1 Unit Test     "}

	// Get JSON
	putStr, _ := json.Marshal(lbPut)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/labels/1", bytes.NewBuffer(putStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.PUT("/api/v3/:account/labels/:id", c.UpdateLabel)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Label{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, result.Name, "Label #1 Unit Test")

	// Double check the db.
	lb := models.Label{}
	db.First(&lb, 1)
	st.Expect(t, lb.Id, uint(1))
	st.Expect(t, lb.Name, "Label #1 Unit Test")
}

//
// TestDeleteLabel01 - Test delete Label 01
//
func TestDeleteLabel01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Test labels. -- First 2 are to make sure we don't get them as they are not our account.
	db.Save(&models.Label{AccountId: 34, Name: "No #1"})
	db.Save(&models.Label{AccountId: 34, Name: "No #2"})
	db.Save(&models.Label{AccountId: 33, Name: "label #1"})
	db.Save(&models.Label{AccountId: 33, Name: "label #2"})
	db.Save(&models.Label{AccountId: 33, Name: "label #3"})
	db.Save(&models.Label{AccountId: 33, Name: "Abc"})
	db.Save(&models.Label{AccountId: 33, Name: "Xyz"})

	// Assign some fake ledger entries to the loop up table.
	db.Save(&models.LabelsToLedger{LabelsToLedgerLabelId: 4, LabelsToLedgerLedgerId: 1})
	db.Save(&models.LabelsToLedger{LabelsToLedgerLabelId: 4, LabelsToLedgerLedgerId: 2})
	db.Save(&models.LabelsToLedger{LabelsToLedgerLabelId: 4, LabelsToLedgerLedgerId: 3})
	db.Save(&models.LabelsToLedger{LabelsToLedgerLabelId: 4, LabelsToLedgerLedgerId: 4})
	db.Save(&models.LabelsToLedger{LabelsToLedgerLabelId: 1, LabelsToLedgerLedgerId: 4})
	db.Save(&models.LabelsToLedger{LabelsToLedgerLabelId: 1, LabelsToLedgerLedgerId: 4})

	// Setup request
	req, _ := http.NewRequest("DELETE", "/api/v3/33/labels/4", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.DELETE("/api/v3/:account/labels/:id", c.DeleteLabel)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 204)

	// Check results
	lbs := []models.Label{}
	db.Find(&lbs)
	st.Expect(t, len(lbs), 6)

	lb := models.Label{}
	db.Where("LabelsName = ?", "label #2").Find(&lb)
	st.Expect(t, lb.Id, uint(0))

	lb2 := models.Label{}
	db.Where("LabelsName = ?", "label #3").Find(&lb2)
	st.Expect(t, lb2.Id, uint(5))

	// Check to make sure LabelsToLedger get deleted
	ltls := []models.LabelsToLedger{}
	db.Find(&ltls)
	st.Expect(t, len(ltls), 2)

	ltl := []models.LabelsToLedger{}
	db.Where("LabelsToLedgerLabelId = ?", 4).Find(&ltl)
	st.Expect(t, len(ltl), 0)
}

//
// TestDeleteLabel02 - Test delete Label 02 - not owned by you
//
func TestDeleteLabel02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Test labels. -- First 2 are to make sure we don't get them as they are not our account.
	db.Save(&models.Label{AccountId: 34, Name: "No #1"})
	db.Save(&models.Label{AccountId: 34, Name: "No #2"})
	db.Save(&models.Label{AccountId: 33, Name: "label #1"})
	db.Save(&models.Label{AccountId: 33, Name: "label #2"})
	db.Save(&models.Label{AccountId: 33, Name: "label #3"})
	db.Save(&models.Label{AccountId: 33, Name: "Abc"})
	db.Save(&models.Label{AccountId: 33, Name: "Xyz"})

	// Assign some fake ledger entries to the loop up table.
	db.Save(&models.LabelsToLedger{LabelsToLedgerLabelId: 4, LabelsToLedgerLedgerId: 1})
	db.Save(&models.LabelsToLedger{LabelsToLedgerLabelId: 4, LabelsToLedgerLedgerId: 2})
	db.Save(&models.LabelsToLedger{LabelsToLedgerLabelId: 4, LabelsToLedgerLedgerId: 3})
	db.Save(&models.LabelsToLedger{LabelsToLedgerLabelId: 4, LabelsToLedgerLedgerId: 4})
	db.Save(&models.LabelsToLedger{LabelsToLedgerLabelId: 1, LabelsToLedgerLedgerId: 4})
	db.Save(&models.LabelsToLedger{LabelsToLedgerLabelId: 1, LabelsToLedgerLedgerId: 4})

	// Setup request
	req, _ := http.NewRequest("DELETE", "/api/v3/33/labels/2", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.DELETE("/api/v3/:account/labels/:id", c.DeleteLabel)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"Label not found."}`)
}

/* End File */
