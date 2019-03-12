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

	"github.com/cloudmanic/skyclerk.com/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
	"github.com/tidwall/gjson"
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
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.GET("/api/v1/:account/categories", c.GetCategories)
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

	// Test categories. -- First 2 are to make sure we don't get them as they are not our account.
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
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.GET("/api/v1/:account/categories", c.GetCategories)
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
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.GET("/api/v1/:account/categories", c.GetCategories)
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
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.GET("/api/v1/:account/categories", c.GetCategories)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Body.String(), `{"error":"There was an error. Please contact help@skyclerk.com for help."}`)
}

//
// Test get a Get Categories 05 - Test Type
//
func TestGetCategories05(t *testing.T) {

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
	req, _ := http.NewRequest("GET", "/api/v1/33/categories?type=income", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.GET("/api/v1/:account/categories", c.GetCategories)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []models.Category{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)

	st.Expect(t, results[0].Id, uint(6))
	st.Expect(t, results[1].Id, uint(4))

	st.Expect(t, results[0].Name, "Abc")
	st.Expect(t, results[1].Name, "Category #2")

	st.Expect(t, results[0].Type, "income")
	st.Expect(t, results[1].Type, "income")
}

//
// Test Get Category 01
//
func TestGetCategory01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test cat to conflict with
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #2"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #3"})

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v1/33/categories/2", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.GET("/api/v1/:account/categories/:id", c.GetCategory)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Category{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, result.Name, "Category #2")
	st.Expect(t, result.Type, "expense")
}

//
// Test Get Category 02 - wrong account
//
func TestGetCategory02(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test cat to conflict with
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1"})
	db.Save(&models.Category{AccountId: 55, Type: "1", Name: "Category #2"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #3"})

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v1/33/categories/2", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.GET("/api/v1/:account/categories/:id", c.GetCategory)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "error").String(), "Category not found.")
}

//
// Test create Category 01
//
func TestCreateCategory01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Post data
	catPost := models.Category{Type: "1", Name: "Category #1"}

	// Get JSON
	postStr, _ := json.Marshal(catPost)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v1/33/categories", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v1/:account/categories", c.CreateCategory)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Category{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 201)
	st.Expect(t, result.Name, "Category #1")
	st.Expect(t, result.Type, "expense")

	// Double check the db.
	cat := models.Category{}
	db.First(&cat, 1)
	st.Expect(t, cat.Id, uint(1))
	st.Expect(t, cat.Name, "Category #1")
	st.Expect(t, cat.Type, "1")
}

//
// Test create Category 02
//
func TestCreateCategory02(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Post data
	catPost := models.Category{Type: "2", Name: "Category #2"}

	// Get JSON
	postStr, _ := json.Marshal(catPost)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v1/33/categories", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v1/:account/categories", c.CreateCategory)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Category{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 201)
	st.Expect(t, result.Name, "Category #2")
	st.Expect(t, result.Type, "income")

	// Double check the db.
	cat := models.Category{}
	db.First(&cat, 1)
	st.Expect(t, cat.Id, uint(1))
	st.Expect(t, cat.Name, "Category #2")
	st.Expect(t, cat.Type, "2")
}

//
// Test create Category 03 - Duplicate category name
//
func TestCreateCategory03(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test cat to conflict with
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1"})

	// Post data
	catPost := models.Category{AccountId: 88, Type: "1", Name: "Category #1"}

	// Get JSON
	postStr, _ := json.Marshal(catPost)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v1/33/categories", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v1/:account/categories", c.CreateCategory)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "errors.name").String(), "Category name is already in use.")
}

//
// Test create Category 04 - Duplicate category name case
//
func TestCreateCategory04(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test cat to conflict with
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1"})

	// Post data
	catPost := models.Category{Type: "1", Name: "category #1"} // lower case c

	// Get JSON
	postStr, _ := json.Marshal(catPost)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v1/33/categories", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v1/:account/categories", c.CreateCategory)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "errors.name").String(), "Category name is already in use.")
}

//
// Test create Category 05 - Duplicate category spaces
//
func TestCreateCategory05(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test cat to conflict with
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1"})

	// Post data
	catPost := models.Category{Type: "1", Name: "  Category #1  "} // spaces

	// Get JSON
	postStr, _ := json.Marshal(catPost)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v1/33/categories", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v1/:account/categories", c.CreateCategory)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "errors.name").String(), "Category name is already in use.")
}

//
// Test create Category 06 - correct types
//
func TestCreateCategory06(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Post data
	catPost := models.Category{Type: "9", Name: "Category #1"}

	// Get JSON
	postStr, _ := json.Marshal(catPost)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v1/33/categories", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v1/:account/categories", c.CreateCategory)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "errors.type").String(), "The type field must be 1, or 2.")
}

//
// Test create Category 07 - No Type
//
func TestCreateCategory07(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Post data
	catPost := models.Category{Name: "Category #1"} // No Type

	// Get JSON
	postStr, _ := json.Marshal(catPost)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v1/33/categories", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v1/:account/categories", c.CreateCategory)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "errors.type").String(), "The type field is required.")
}

//
// Test create Category 08 - Type with spaces
//
func TestCreateCategory08(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Post data
	catPost := models.Category{Type: " 1 ", Name: "Category #1"} // spaces

	// Get JSON
	postStr, _ := json.Marshal(catPost)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v1/33/categories", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v1/:account/categories", c.CreateCategory)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "errors.type").String(), "The type field must be 1, or 2.")
}

//
// Test create Category 09 - Same cat, different type
//
func TestCreateCategory09(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test cat to conflict with
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1"})

	// Post data
	catPost := models.Category{Type: "2", Name: "Category #1"}

	// Get JSON
	postStr, _ := json.Marshal(catPost)

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v1/33/categories", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.POST("/api/v1/:account/categories", c.CreateCategory)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Category{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 201)
	st.Expect(t, result.Name, "Category #1")
	st.Expect(t, result.Type, "income")
}

//
// Test update Category 01
//
func TestUpdateCategory01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test cat to conflict with
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1"})

	// Put data
	catPut := models.Category{Type: "1", Name: "Category #1 Unit Test"}

	// Get JSON
	putStr, _ := json.Marshal(catPut)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v1/33/categories/1", bytes.NewBuffer(putStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.PUT("/api/v1/:account/categories/:id", c.UpdateCategory)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Category{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, result.Name, "Category #1 Unit Test")
	st.Expect(t, result.Type, "expense")

	// Double check the db.
	cat := models.Category{}
	db.First(&cat, 1)
	st.Expect(t, cat.Id, uint(1))
	st.Expect(t, cat.Name, "Category #1 Unit Test")
	st.Expect(t, cat.Type, "1")
}

//
// Test update Category 02 - Duplicate category
//
func TestUpdateCategory02(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test cat to conflict with
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1 Unit Test"})

	// Put data
	catPut := models.Category{Type: "1", Name: "Category #1 Unit Test"}

	// Get JSON
	putStr, _ := json.Marshal(catPut)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v1/33/categories/1", bytes.NewBuffer(putStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.PUT("/api/v1/:account/categories/:id", c.UpdateCategory)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "errors.name").String(), "Category name is already in use.")
}

//
// Test update Category 03 - Duplicate category same id
//
func TestUpdateCategory03(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test cat to conflict with
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1 Unit Test"})

	// Put data
	catPut := models.Category{Type: "1", Name: "Category #1 Unit Test"}

	// Get JSON
	putStr, _ := json.Marshal(catPut)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v1/33/categories/2", bytes.NewBuffer(putStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.PUT("/api/v1/:account/categories/:id", c.UpdateCategory)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Category{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, result.Id, uint(2))
	st.Expect(t, result.Name, "Category #1 Unit Test")
	st.Expect(t, result.Type, "expense")

	// Double check the db.
	cat := models.Category{}
	db.First(&cat, 2)
	st.Expect(t, cat.Id, uint(2))
	st.Expect(t, cat.Name, "Category #1 Unit Test")
	st.Expect(t, cat.Type, "1")
}

//
// Test update Category 04 - Duplicate category same id new type
//
func TestUpdateCategory04(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test cat to conflict with
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1 Unit Test"})

	// Put data
	catPut := models.Category{Type: "2", Name: "Category #1 Unit Test"}

	// Get JSON
	putStr, _ := json.Marshal(catPut)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v1/33/categories/2", bytes.NewBuffer(putStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.PUT("/api/v1/:account/categories/:id", c.UpdateCategory)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.Category{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, w.Code, 200)
	st.Expect(t, result.Id, uint(2))
	st.Expect(t, result.Name, "Category #1 Unit Test")
	st.Expect(t, result.Type, "income")

	// Double check the db.
	cat := models.Category{}
	db.First(&cat, 2)
	st.Expect(t, cat.Id, uint(2))
	st.Expect(t, cat.Name, "Category #1 Unit Test")
	st.Expect(t, cat.Type, "2")
}

//
// Test update Category 05 - Duplicate category spaces
//
func TestUpdateCategory05(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test cat to conflict with
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1 Unit Test"})

	// Put data
	catPut := models.Category{Type: "1", Name: "    Category #1 Unit Test     "} // spaces

	// Get JSON
	putStr, _ := json.Marshal(catPut)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v1/33/categories/1", bytes.NewBuffer(putStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.PUT("/api/v1/:account/categories/:id", c.UpdateCategory)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "errors.name").String(), "Category name is already in use.")
}

//
// Test update Category 06 - Duplicate category casing
//
func TestUpdateCategory06(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test cat to conflict with
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1 Unit Test"})

	// Put data
	catPut := models.Category{Type: "1", Name: "category #1 unit test"} // casing

	// Get JSON
	putStr, _ := json.Marshal(catPut)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v1/33/categories/1", bytes.NewBuffer(putStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.PUT("/api/v1/:account/categories/:id", c.UpdateCategory)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "errors.name").String(), "Category name is already in use.")
}

//
// Test delete Category 01
//
func TestDeleteCategory01(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test cat to conflict with
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #2 Delete Me"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #3"})

	// Setup request
	req, _ := http.NewRequest("DELETE", "/api/v1/33/categories/2", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.DELETE("/api/v1/:account/categories/:id", c.DeleteCategory)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 204)

	// Double check the db.
	cats := []models.Category{}
	db.Find(&cats)
	st.Expect(t, len(cats), 2)
	st.Expect(t, cats[0].Id, uint(1))
	st.Expect(t, cats[0].Name, "Category #1")
	st.Expect(t, cats[1].Id, uint(3))
	st.Expect(t, cats[1].Name, "Category #3")
}

//
// Test delete Category 02
//
func TestDeleteCategory02(t *testing.T) {

	// Start the db connection.
	db, _ := models.NewDB()
	defer db.Close()

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test cat to conflict with
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #1"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #2 Delete Me"})
	db.Save(&models.Category{AccountId: 33, Type: "1", Name: "Category #3"})

	// Create test ledger
	db.Save(&models.Ledger{AccountId: 33, CategoryId: 2})

	// Setup request
	req, _ := http.NewRequest("DELETE", "/api/v1/33/categories/2", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", 109)
	})
	r.DELETE("/api/v1/:account/categories/:id", c.DeleteCategory)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, gjson.Get(w.Body.String(), "error").String(), "Can not delete category. It is in use by a ledger entry.")
}

/* End File */
