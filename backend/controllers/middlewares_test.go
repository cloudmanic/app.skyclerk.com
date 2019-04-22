//
// Date: 3/29/2019
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"app.skyclerk.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// TestAuthMiddleware01 - Success login
//
func TestAuthMiddleware01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create a test user in DB (this function is in oauth_test.go)
	user, account, app := createTestUserInDB(db)

	// Create a test session
	sess, err := db.CreateSession(user.Id, app.Id, "Test UserAgent", "1.2.3.4")
	st.Expect(t, err, nil)

	// Setup request
	req, _ := http.NewRequest("GET", fmt.Sprintf("/%d/ping", account.Id), nil)
	req.Header.Set("Authorization", "Bearer "+sess.AccessToken)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(c.AuthMiddleware())
	r.GET("/:account/ping", c.PingFromServer)
	r.ServeHTTP(w, req)

	// Validate results
	st.Expect(t, w.Body.String(), `{"status":"ok"}`)
}

//
// TestAuthMiddleware02 - Failed login
//
func TestAuthMiddleware02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create a test user in DB (this function is in oauth_test.go)
	user, account, app := createTestUserInDB(db)

	// Create a test session
	_, err := db.CreateSession(user.Id, app.Id, "Test UserAgent", "1.2.3.4")
	st.Expect(t, err, nil)

	// Setup request
	req, _ := http.NewRequest("GET", fmt.Sprintf("/%d/ping", account.Id), nil)
	req.Header.Set("Authorization", "Bearer blahblahblah")

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(c.AuthMiddleware())
	r.GET("/:account/ping", c.PingFromServer)
	r.ServeHTTP(w, req)

	// Validate results
	st.Expect(t, w.Body.String(), `{"errors":{"system":"Authorization Failed (#003)"}}`)
}

//
// TestAuthMiddleware03 - Failed account
//
func TestAuthMiddleware03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create a test user in DB (this function is in oauth_test.go)
	user, _, app := createTestUserInDB(db)

	// Create a test session
	sess, err := db.CreateSession(user.Id, app.Id, "Test UserAgent", "1.2.3.4")
	st.Expect(t, err, nil)

	// Setup request
	req, _ := http.NewRequest("GET", fmt.Sprintf("/%d/ping", 99), nil)
	req.Header.Set("Authorization", "Bearer "+sess.AccessToken)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(c.AuthMiddleware())
	r.GET("/:account/ping", c.PingFromServer)
	r.ServeHTTP(w, req)

	// Validate results
	st.Expect(t, w.Body.String(), `{"errors":{"system":"Account Not Found - Unable to Authenticate (#006)"}}`)
}

//
// TestAuthMiddleware04 - No access token
//
func TestAuthMiddleware04(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create a test user in DB (this function is in oauth_test.go)
	user, _, app := createTestUserInDB(db)

	// Create a test session
	_, err := db.CreateSession(user.Id, app.Id, "Test UserAgent", "1.2.3.4")
	st.Expect(t, err, nil)

	// Setup request - missing access token
	req, _ := http.NewRequest("GET", fmt.Sprintf("/%d/ping", 99), nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(c.AuthMiddleware())
	r.GET("/:account/ping", c.PingFromServer)
	r.ServeHTTP(w, req)

	// Validate results
	st.Expect(t, w.Body.String(), `{"errors":{"system":"Authorization Failed (#001)"}}`)
}

//
// TestAuthNoAccountMiddleware01 - Success login
//
func TestAuthNoAccountMiddleware01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create a test user in DB (this function is in oauth_test.go)
	user, account, app := createTestUserInDB(db)

	// Create a test session
	sess, err := db.CreateSession(user.Id, app.Id, "Test UserAgent", "1.2.3.4")
	st.Expect(t, err, nil)

	// Setup request
	req, _ := http.NewRequest("GET", fmt.Sprintf("/%d/ping", account.Id), nil)
	req.Header.Set("Authorization", "Bearer "+sess.AccessToken)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(c.AuthNoAccountMiddleware())
	r.GET("/:account/ping", c.PingFromServer)
	r.ServeHTTP(w, req)

	// Validate results
	st.Expect(t, w.Body.String(), `{"status":"ok"}`)
}

//
// TestAuthNoAccountMiddleware02 - Failed login
//
func TestAuthNoAccountMiddleware02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create a test user in DB (this function is in oauth_test.go)
	user, account, app := createTestUserInDB(db)

	// Create a test session
	_, err := db.CreateSession(user.Id, app.Id, "Test UserAgent", "1.2.3.4")
	st.Expect(t, err, nil)

	// Setup request
	req, _ := http.NewRequest("GET", fmt.Sprintf("/%d/ping", account.Id), nil)
	req.Header.Set("Authorization", "Bearer blahblahblah")

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(c.AuthNoAccountMiddleware())
	r.GET("/:account/ping", c.PingFromServer)
	r.ServeHTTP(w, req)

	// Validate results
	st.Expect(t, w.Body.String(), `{"errors":{"system":"Authorization Failed (#003)"}}`)
}

//
// TestAuthNoAccountMiddleware03 - Failed account
//
func TestAuthNoAccountMiddleware03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create a test user in DB (this function is in oauth_test.go)
	user, _, app := createTestUserInDB(db)

	// Create a test session
	_, err := db.CreateSession(user.Id, app.Id, "Test UserAgent", "1.2.3.4")
	st.Expect(t, err, nil)

	// Setup request - missing access token
	req, _ := http.NewRequest("GET", fmt.Sprintf("/%d/ping", 99), nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(c.AuthNoAccountMiddleware())
	r.GET("/:account/ping", c.PingFromServer)
	r.ServeHTTP(w, req)

	// Validate results
	st.Expect(t, w.Body.String(), `{"errors":{"system":"Authorization Failed (#001)"}}`)
}

/* End File */
