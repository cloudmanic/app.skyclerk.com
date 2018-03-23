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
// Test get a Get Categories 01
//
func TestGetCategories01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Test labels. -- First 2 are to make sure we don't get them as they are not our account.
	db.Save(&models.Category{AccountId: 34, Type: "1", Name: "No #1"})
	db.Save(&models.Category{AccountId: 34, Type: "1", Name: "No #2"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1"})
	db.Save(&models.Category{AccountId: 33, Type: "2", Name: "Category #2"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #3"})
	db.Save(&models.Category{AccountId: 33, Type: "2", Name: "Abc"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Xyz"})

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v1/33/categories", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("account", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v1/33/categories", c.GetCategories)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []models.Category{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)

	st.Expect(t, results[0].Id, uint(6))
	st.Expect(t, results[1].Id, uint(3))
	st.Expect(t, results[2].Id, uint(4))
	st.Expect(t, results[3].Id, uint(5))
	st.Expect(t, results[4].Id, uint(7))

	st.Expect(t, results[0].Name, "Abc")
	st.Expect(t, results[1].Name, "Category #1")
	st.Expect(t, results[2].Name, "Category #2")
	st.Expect(t, results[3].Name, "Category #3")
	st.Expect(t, results[4].Name, "Xyz")

	st.Expect(t, results[0].Type, "income")
	st.Expect(t, results[1].Type, "expense")
	st.Expect(t, results[2].Type, "income")
	st.Expect(t, results[3].Type, "expense")
	st.Expect(t, results[4].Type, "expense")
}

//
// Test get a Get Categories 02 - Test Sort
//
func TestGetCategories02(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Test labels. -- First 2 are to make sure we don't get them as they are not our account.
	db.Save(&models.Category{AccountId: 34, Type: "1", Name: "No #1"})
	db.Save(&models.Category{AccountId: 34, Type: "1", Name: "No #2"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1"})
	db.Save(&models.Category{AccountId: 33, Type: "2", Name: "Category #2"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #3"})
	db.Save(&models.Category{AccountId: 33, Type: "2", Name: "Abc"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Xyz"})

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v1/33/categories?sort=desc", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("account", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v1/33/categories", c.GetCategories)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []models.Category{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)

	st.Expect(t, results[0].Id, uint(7))
	st.Expect(t, results[1].Id, uint(5))
	st.Expect(t, results[2].Id, uint(4))
	st.Expect(t, results[3].Id, uint(3))
	st.Expect(t, results[4].Id, uint(6))

	st.Expect(t, results[0].Name, "Xyz")
	st.Expect(t, results[1].Name, "Category #3")
	st.Expect(t, results[2].Name, "Category #2")
	st.Expect(t, results[3].Name, "Category #1")
	st.Expect(t, results[4].Name, "Abc")

	st.Expect(t, results[0].Type, "expense")
	st.Expect(t, results[1].Type, "expense")
	st.Expect(t, results[2].Type, "income")
	st.Expect(t, results[3].Type, "expense")
	st.Expect(t, results[4].Type, "income")
}

//
// Test get a Get Categories 03 - Test Sort
//
func TestGetCategories03(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Test labels. -- First 2 are to make sure we don't get them as they are not our account.
	db.Save(&models.Category{AccountId: 34, Type: "1", Name: "No #1"})
	db.Save(&models.Category{AccountId: 34, Type: "1", Name: "No #2"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1"})
	db.Save(&models.Category{AccountId: 33, Type: "2", Name: "Category #2"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #3"})
	db.Save(&models.Category{AccountId: 33, Type: "2", Name: "Abc"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Xyz"})

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v1/33/categories?sort=desc&order=CategoriesId", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("account", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v1/33/categories", c.GetCategories)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []models.Category{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)

	st.Expect(t, results[0].Id, uint(7))
	st.Expect(t, results[1].Id, uint(6))
	st.Expect(t, results[2].Id, uint(5))
	st.Expect(t, results[3].Id, uint(4))
	st.Expect(t, results[4].Id, uint(3))

	st.Expect(t, results[0].Name, "Xyz")
	st.Expect(t, results[1].Name, "Abc")
	st.Expect(t, results[2].Name, "Category #3")
	st.Expect(t, results[3].Name, "Category #2")
	st.Expect(t, results[4].Name, "Category #1")
}

//
// Test get a Get Categories 04 - Failed Order Col
//
func TestGetCategories04(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Test labels. -- First 2 are to make sure we don't get them as they are not our account.
	db.Save(&models.Category{AccountId: 34, Type: "1", Name: "No #1"})
	db.Save(&models.Category{AccountId: 34, Type: "1", Name: "No #2"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1"})
	db.Save(&models.Category{AccountId: 33, Type: "2", Name: "Category #2"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #3"})
	db.Save(&models.Category{AccountId: 33, Type: "2", Name: "Abc"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Xyz"})

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v1/33/categories?sort=desc&order=FailedId", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("account", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v1/33/categories", c.GetCategories)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Body.String(), `{"error":"There was an error. Please contact help@skyclerk.com for help."}`)
}

/* End File */
