//
// Date: 6/23/2018
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2018 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"app.skyclerk.com/backend/library/test"
	"app.skyclerk.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
	"github.com/tidwall/gjson"
)

//
// Create a test AcctUser
//
func createTestUserInDB(db *models.DB) (models.User, models.Account, models.Application) {
	// Setup test data
	user := test.GetRandomUser(33)
	db.Save(&user)

	account := test.GetRandomAccount(33)
	account.OwnerId = user.Id
	db.Save(&account)

	app := test.GetRandomApplication()
	db.Save(&app)

	// Set new user account to users
	lu := models.AcctToUsers{AcctId: account.Id, UserId: user.Id}
	db.Save(&lu)

	return user, account, app
}

//
// TestDoOauthToken01 test to make sure auth works. - Success login
//
func TestDoOauthToken01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create a test user in DB
	user, account, app := createTestUserInDB(db)

	// Build stuct that we convert to json.
	type PostStruct struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		GrantType string `json:"grant_type"`
		ClientID  string `json:"client_id"`
	}

	pS := PostStruct{
		Username:  user.Email,
		Password:  "F00bAr123",
		GrantType: "password",
		ClientID:  app.ClientId,
	}

	// Get JSON
	postStr, _ := json.Marshal(pS)

	// Setup request
	req, _ := http.NewRequest("POST", "/oauth/token", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.POST("/oauth/token", c.DoOauthToken)
	r.ServeHTTP(w, req)

	// Get results from json.
	userId := gjson.Get(w.Body.String(), "user_id").Int()
	tokenType := gjson.Get(w.Body.String(), "token_type").String()
	accessToken := gjson.Get(w.Body.String(), "access_token").String()

	// Test results
	st.Expect(t, w.Code, 200)
	st.Expect(t, tokenType, "bearer")
	st.Expect(t, uint(userId), user.Id)
	st.Expect(t, len(accessToken), 50)

	// Look in the sessions table for the correct access token
	sess, err := db.GetByAccessToken(accessToken)
	st.Expect(t, err, nil)
	st.Expect(t, sess.UserId, user.Id)
	st.Expect(t, sess.ApplicationId, app.Id)
	st.Expect(t, sess.AccessToken, accessToken)

	// Let's get the user to make sure the accounts are added
	u, err := db.GetUserById(user.Id)
	st.Expect(t, err, nil)
	st.Expect(t, len(u.Accounts), 1)
	st.Expect(t, u.Accounts[0].Name, account.Name)
}

//
// TestDoOauthToken02 test to make sure auth works. - Failed login.
//
func TestDoOauthToken02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create a test user in DB
	user, _, app := createTestUserInDB(db)

	// Build stuct that we convert to json.
	type PostStruct struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		GrantType string `json:"grant_type"`
		ClientID  string `json:"client_id"`
	}

	pS := PostStruct{
		Username:  user.Email,
		Password:  "nogood",
		GrantType: "password",
		ClientID:  app.ClientId,
	}

	// Get JSON
	postStr, _ := json.Marshal(pS)

	// Setup request
	req, _ := http.NewRequest("POST", "/oauth/token", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.POST("/oauth/token", c.DoOauthToken)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"Sorry, we could not find your account."}`)
}

//
// TestDoOauthToken03 test to make sure auth works. - Failed login (small password)
//
func TestDoOauthToken03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create a test user in DB
	user, _, app := createTestUserInDB(db)

	// Build stuct that we convert to json.
	type PostStruct struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		GrantType string `json:"grant_type"`
		ClientID  string `json:"client_id"`
	}

	pS := PostStruct{
		Username:  user.Email,
		Password:  "a",
		GrantType: "password",
		ClientID:  app.ClientId,
	}

	// Get JSON
	postStr, _ := json.Marshal(pS)

	// Setup request
	req, _ := http.NewRequest("POST", "/oauth/token", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.POST("/oauth/token", c.DoOauthToken)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"The password filed must be at least 6 characters long."}`)
}

//
// TestDoOauthToken04 test to make sure auth works. - Failed login (bad email)
//
func TestDoOauthToken04(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create a test user in DB
	_, _, app := createTestUserInDB(db)

	// Build stuct that we convert to json.
	type PostStruct struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		GrantType string `json:"grant_type"`
		ClientID  string `json:"client_id"`
	}

	pS := PostStruct{
		Username:  "bademail@example.com",
		Password:  "F00bAr123",
		GrantType: "password",
		ClientID:  app.ClientId,
	}

	// Get JSON
	postStr, _ := json.Marshal(pS)

	// Setup request
	req, _ := http.NewRequest("POST", "/oauth/token", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.POST("/oauth/token", c.DoOauthToken)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"Sorry, we were unable to find our account."}`)
}

//
// TestDoOauthToken05 test to make sure auth works. - Failed login (bad grant type)
//
func TestDoOauthToken05(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create a test user in DB
	user, _, app := createTestUserInDB(db)

	// Build stuct that we convert to json.
	type PostStruct struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		GrantType string `json:"grant_type"`
		ClientID  string `json:"client_id"`
	}

	pS := PostStruct{
		Username:  user.Email,
		Password:  "F00bAr123",
		GrantType: "token",
		ClientID:  app.ClientId,
	}

	// Get JSON
	postStr, _ := json.Marshal(pS)

	// Setup request
	req, _ := http.NewRequest("POST", "/oauth/token", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.POST("/oauth/token", c.DoOauthToken)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"Invalid client_id or grant type."}`)
}

//
// TestDoOauthToken06 test to make sure auth works. - Success login (more than one account)
//
func TestDoOauthToken06(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create a test user in DB
	user, account, app := createTestUserInDB(db)

	// Add more accounts.
	a1 := test.GetRandomAccount(20)
	db.Save(&a1)
	a2 := test.GetRandomAccount(21)
	db.Save(&a2)
	a3 := test.GetRandomAccount(23)
	db.Save(&a3)

	lu1 := models.AcctToUsers{AcctId: a1.Id, UserId: user.Id}
	db.Save(&lu1)

	lu2 := models.AcctToUsers{AcctId: a2.Id, UserId: user.Id}
	db.Save(&lu2)

	lu3 := models.AcctToUsers{AcctId: a3.Id, UserId: uint(5)}
	db.Save(&lu3)

	// Build stuct that we convert to json.
	type PostStruct struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		GrantType string `json:"grant_type"`
		ClientID  string `json:"client_id"`
	}

	pS := PostStruct{
		Username:  user.Email,
		Password:  "F00bAr123",
		GrantType: "password",
		ClientID:  app.ClientId,
	}

	// Get JSON
	postStr, _ := json.Marshal(pS)

	// Setup request
	req, _ := http.NewRequest("POST", "/oauth/token", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.POST("/oauth/token", c.DoOauthToken)
	r.ServeHTTP(w, req)

	// Get results from json.
	userId := gjson.Get(w.Body.String(), "user_id").Int()
	tokenType := gjson.Get(w.Body.String(), "token_type").String()
	accessToken := gjson.Get(w.Body.String(), "access_token").String()

	// Test results
	st.Expect(t, w.Code, 200)
	st.Expect(t, tokenType, "bearer")
	st.Expect(t, uint(userId), user.Id)
	st.Expect(t, len(accessToken), 50)

	// Look in the sessions table for the correct access token
	sess, err := db.GetByAccessToken(accessToken)
	st.Expect(t, err, nil)
	st.Expect(t, sess.UserId, user.Id)
	st.Expect(t, sess.ApplicationId, app.Id)
	st.Expect(t, sess.AccessToken, accessToken)

	// Let's get the user to make sure the accounts are added
	u, err := db.GetUserById(user.Id)
	st.Expect(t, err, nil)
	st.Expect(t, len(u.Accounts), 3)
	st.Expect(t, u.Accounts[0].Name, account.Name)
	st.Expect(t, u.Accounts[1].Name, a2.Name)
	st.Expect(t, u.Accounts[2].Name, a1.Name)
}

//
// TestDoLogOut01 - test logout
//
func TestDoLogOut01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create a test user in DB
	user, _, app := createTestUserInDB(db)

	// Build stuct that we convert to json.
	type PostStruct struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		GrantType string `json:"grant_type"`
		ClientID  string `json:"client_id"`
	}

	pS := PostStruct{
		Username:  user.Email,
		Password:  "F00bAr123",
		GrantType: "password",
		ClientID:  app.ClientId,
	}

	// Get JSON
	postStr, _ := json.Marshal(pS)

	// Setup request
	req, _ := http.NewRequest("POST", "/oauth/token", bytes.NewBuffer(postStr))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.POST("/oauth/token", c.DoOauthToken)
	r.ServeHTTP(w, req)

	// Get results from json.
	userId := gjson.Get(w.Body.String(), "user_id").Int()
	accessToken := gjson.Get(w.Body.String(), "access_token").String()

	// Test results
	st.Expect(t, w.Code, 200)
	st.Expect(t, uint(userId), user.Id)
	st.Expect(t, len(accessToken), 50)

	// Setup request
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/oauth/logout?access_token="+accessToken, nil)
	r.GET("/oauth/logout", c.DoLogOut)
	r.ServeHTTP(w2, req2)

	// Test results
	st.Expect(t, w2.Code, 200)
	st.Expect(t, w2.Body.String(), `{"status":"ok"}`)

	// Look in the sessions table for the correct access token
	sess, err := db.GetByAccessToken(accessToken)
	st.Expect(t, err.Error(), "Access Token Not Found - Unable to Authenticate (#001)")
	st.Expect(t, sess.UserId, uint(0))
	st.Expect(t, sess.AccessToken, "")

	// Setup request
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/oauth/logout?access_token="+accessToken, nil)
	r.ServeHTTP(w3, req3)

	// Test results
	st.Expect(t, w3.Code, 400)
	st.Expect(t, w3.Body.String(), `{"error":"Sorry, we could not find your session.","status":"error"}`)

	// Setup request
	w4 := httptest.NewRecorder()
	req4, _ := http.NewRequest("GET", "/oauth/logout?access_token=", nil)
	r.ServeHTTP(w4, req4)

	// Test results
	st.Expect(t, w4.Code, 400)
	st.Expect(t, w4.Body.String(), `{"error":"Sorry, access_token is required."}`)
}

/* End File */
