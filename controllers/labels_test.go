//
// Date: 2018-03-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-22
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudmanic/skyclerk.com/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// Test get a ledger 01
//
func TestGetLabels01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

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
	req, _ := http.NewRequest("GET", "/api/v1/33/labels", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("account", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v1/33/labels", c.GetLabels)
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
	db, _ := models.NewDB()
	defer db.Close()

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
	req, _ := http.NewRequest("GET", "/api/v1/33/labels?sort=desc", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("account", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v1/33/labels", c.GetLabels)
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
	db, _ := models.NewDB()
	defer db.Close()

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
	req, _ := http.NewRequest("GET", "/api/v1/33/labels?sort=desc&order=LabelsId", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("account", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v1/33/labels", c.GetLabels)
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
	db, _ := models.NewDB()
	defer db.Close()

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
	req, _ := http.NewRequest("GET", "/api/v1/33/labels?sort=asc&order=LabelsId", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("account", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v1/33/labels", c.GetLabels)
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
	db, _ := models.NewDB()
	defer db.Close()

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
	req, _ := http.NewRequest("GET", "/api/v1/33/labels?sort=asc&order=FailedId", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("account", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v1/33/labels", c.GetLabels)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Body.String(), `{"error":"There was an error. Please contact help@skyclerk.com for help."}`)
}

/* End File */
