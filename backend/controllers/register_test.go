//
// Date: 2019-09-13
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
	"golang.org/x/crypto/bcrypt"
)

//
// TestDoRegister01 Test registring a new user.
//
func TestDoRegister01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create applicaiton.
	app := test.GetRandomApplication()
	app.GrantType = "password"
	db.Save(&app)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Get JSON
	postStr := fmt.Sprintf(`{ "first": "%s", "last": "%s", "email": "%s", "password": "%s", "client_id": "%s" }`, "Jane", "Wells", "jane@wells.com", "foobar123", app.ClientId)

	// Setup request
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.POST("/register", c.DoRegister)
	r.ServeHTTP(w, req)

	// Decode response.
	type Response struct {
		UserId      uint   `json:"user_id"`
		AccessToken string `json:"access_token"`
		AccountId   uint   `json:"account_id"`
	}
	var res Response
	json.Unmarshal(w.Body.Bytes(), &res)

	// Check the database that proper entries where created
	u := models.AcctToUsers{}
	db.Where("acct_id = ? AND user_id = ?", 1, 1).First(&u)

	// Check the database that proper entries where created
	s := models.Session{}
	db.Where("user_id = ? AND application_id = ?", 1, 1).First(&s)

	// Check the database that proper entries where created
	m := models.User{}
	db.Where("id = ?", 1).First(&m)

	// Check the database that proper entries where created
	a := models.Account{}
	db.Where("owner_id = ?", 1).First(&a)

	// Test results
	st.Expect(t, w.Code, 200)
	st.Expect(t, res.UserId, uint(1))
	st.Expect(t, res.AccountId, uint(1))
	st.Expect(t, res.AccessToken, s.AccessToken)
	st.Expect(t, u.Id, uint(1))
	st.Expect(t, s.Id, uint(1))
	st.Expect(t, m.Id, uint(1))
	st.Expect(t, m.FirstName, "Jane")
	st.Expect(t, m.LastName, "Wells")
	st.Expect(t, m.Email, "jane@wells.com")
	st.Expect(t, a.Name, "Jane's Skyclerk")

	// Test password.
	err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte("foobar123"))
	st.Expect(t, err, nil)
}

//
// TestDoRegister02 - Error 01 (bad email)
//
func TestDoRegister02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create applicaiton.
	app := test.GetRandomApplication()
	app.GrantType = "password"
	db.Save(&app)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Get JSON
	postStr := fmt.Sprintf(`{ "first": "%s", "last": "%s", "email": "%s", "password": "%s", "client_id": "%s" }`, "Jane", "Wells", "jane", "foobar123", app.ClientId)

	// Setup request
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.POST("/register", c.DoRegister)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"Email address is not a valid format."}`)
}

//
// TestDoRegister03 - Error 02 (no first)
//
func TestDoRegister03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create applicaiton.
	app := test.GetRandomApplication()
	app.GrantType = "password"
	db.Save(&app)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Get JSON
	postStr := fmt.Sprintf(`{ "first": "%s", "last": "%s", "email": "%s", "password": "%s", "client_id": "%s" }`, "", "Wells", "jane@wells.com", "foobar123", app.ClientId)

	// Setup request
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.POST("/register", c.DoRegister)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"First name field is required."}`)
}

//
// TestDoRegister04 - Error 03 (bad password)
//
func TestDoRegister04(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create applicaiton.
	app := test.GetRandomApplication()
	app.GrantType = "password"
	db.Save(&app)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Get JSON
	postStr := fmt.Sprintf(`{ "first": "%s", "last": "%s", "email": "%s", "password": "%s", "client_id": "%s" }`, "Jane", "Wells", "jane@wells.com", "ff", app.ClientId)

	// Setup request
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.POST("/register", c.DoRegister)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"The password filed must be at least 6 characters long."}`)
}

//
// TestDoRegister05 - Error 03 (no last)
//
func TestDoRegister05(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create applicaiton.
	app := test.GetRandomApplication()
	app.GrantType = "password"
	db.Save(&app)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Get JSON
	postStr := fmt.Sprintf(`{ "first": "%s", "last": "%s", "email": "%s", "password": "%s", "client_id": "%s"  }`, "Jane", "", "jane@wells.com", "foobar123", app.ClientId)

	// Setup request
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.POST("/register", c.DoRegister)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"Last name field is required."}`)
}

//
// TestDoRegister06 - Error 04 (bad client id)
//
func TestDoRegister06(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create applicaiton.
	app := test.GetRandomApplication()
	app.GrantType = "password"
	db.Save(&app)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Get JSON
	postStr := fmt.Sprintf(`{ "first": "%s", "last": "%s", "email": "%s", "password": "%s", "client_id": "%s"  }`, "Jane", "", "jane@wells.com", "foobar123", "bad")

	// Setup request
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.POST("/register", c.DoRegister)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"Something went wrong while logging into your account. Please try again or contact help@options.cafe. Sorry for the trouble."}`)
}

//
// TestDoRegister07 - Error 05 (missing client id)
//
func TestDoRegister07(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create applicaiton.
	app := test.GetRandomApplication()
	app.GrantType = "password"
	db.Save(&app)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Get JSON
	postStr := fmt.Sprintf(`{ "first": "%s", "last": "%s", "email": "%s", "password": "%s" }`, "Jane", "", "jane@wells.com", "foobar123")

	// Setup request
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.POST("/register", c.DoRegister)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"Something went wrong while logging into your account. Please try again or contact help@options.cafe. Sorry for the trouble."}`)
}

//
// TestDoRegister08 - Error user already in the system.
//
func TestDoRegister08(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create applicaiton.
	app := test.GetRandomApplication()
	app.GrantType = "password"
	db.Save(&app)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Get JSON
	postStr := fmt.Sprintf(`{ "first": "%s", "last": "%s", "email": "%s", "password": "%s", "client_id": "%s" }`, "Jane", "Wells", "jane@wells.com", "foobar123", app.ClientId)

	// Setup request
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.POST("/register", c.DoRegister)
	r.ServeHTTP(w, req)

	// --------- Register again so we get errors ---------- //

	// Setup request
	req1, _ := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w1 := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r1 := gin.New()
	r1.POST("/register", c.DoRegister)
	r1.ServeHTTP(w1, req1)

	// Test results
	st.Expect(t, w1.Code, 400)
	st.Expect(t, w1.Body.String(), `{"error":"Looks like you already have an account."}`)
}

/* End File */
