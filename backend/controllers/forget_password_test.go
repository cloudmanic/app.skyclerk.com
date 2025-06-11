//
// Date: 10/26/2019
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"app.skyclerk.com/backend/library/test"
	"app.skyclerk.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
	"github.com/tidwall/gjson"
)

//
// TestDoForgotPassword01 - test sending in an email and sending a forgot password link.
//
func TestDoForgotPassword01(t *testing.T) {
	// Skip if no mail driver configured for testing
	if len(os.Getenv("MAIL_DRIVER")) == 0 {
		t.Skip("Skipping test - no mail driver configured")
		return
	}

	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test users.
	u1 := test.GetRandomUser(33)
	db.Save(&u1)

	// JSON we post in
	postStr := fmt.Sprintf(`{ "email": "%s" }`, u1.Email)

	// Setup request
	req, _ := http.NewRequest("POST", "/forgot-password", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.POST("/forgot-password", c.DoForgotPassword)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 204)

	// Double check the db.
	l := models.ForgotPassword{}
	db.First(&l, 1)
	st.Expect(t, l.Id, uint(1))
	st.Expect(t, len(l.Token), 30)
}

//
// TestDoForgotPassword02 - Email not found
//
func TestDoForgotPassword02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test users.
	u1 := test.GetRandomUser(33)
	db.Save(&u1)

	// JSON we post in
	postStr := fmt.Sprintf(`{ "email": "%s" }`, "woots@example.com")

	// Setup request
	req, _ := http.NewRequest("POST", "/forgot-password", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.POST("/forgot-password", c.DoForgotPassword)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"Sorry, we could not find your account."}`)
}

//
// TestDoResetPassword01 - Test restsetting the password after we did the email dance.
//
func TestDoResetPassword01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test users.
	u1 := test.GetRandomUser(33)
	db.Save(&u1)

	// Set reset password token
	token := "wootsblahdope"
	y := models.ForgotPassword{Token: token, UserId: u1.Id}
	db.Save(&y)

	// JSON we post in
	postStr := fmt.Sprintf(`{ "hash": "%s", "password": "helloworld123" }`, token)

	// Setup request
	req, _ := http.NewRequest("POST", "/reset-password", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.POST("/reset-password", c.DoResetPassword)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 204)

	// Double check the db.
	l := models.ForgotPassword{}
	db.First(&l, 1)
	st.Expect(t, l.Id, uint(0))
	st.Expect(t, len(l.Token), 0)

	i := models.User{}
	db.First(&i, 1)
	st.Expect(t, i.Id, uint(1))
	st.Expect(t, len(i.Md5Password), 0)
	st.Expect(t, len(i.Md5Salt), 0)
	st.Expect(t, (len(i.Password) > 0), true)

	// ---------- Double check the auth worked. ------------ //

	// Build random app
	app := test.GetRandomApplication()
	db.Save(&app)

	// Build random account
	account := test.GetRandomAccount(33)
	account.OwnerId = u1.Id
	db.Save(&account)

	// Set new user account to users
	lu := models.AcctToUsers{AccountId: account.Id, UserId: u1.Id}
	db.Save(&lu)

	// Build stuct that we convert to json.
	type PostStruct struct {
		Username  string `json:"username"`
		Password  string `json:"password"`
		GrantType string `json:"grant_type"`
		ClientID  string `json:"client_id"`
	}

	pS := PostStruct{
		Username:  u1.Email,
		Password:  "helloworld123",
		GrantType: "password",
		ClientID:  app.ClientId,
	}

	// Get JSON
	postStr2, _ := json.Marshal(pS)

	// Setup request
	req2, _ := http.NewRequest("POST", "/oauth/token", bytes.NewBuffer(postStr2))

	// Setup writer.
	w = httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r = gin.New()
	r.POST("/oauth/token", c.DoOauthToken)
	r.ServeHTTP(w, req2)

	// Get results from json.
	userId := gjson.Get(w.Body.String(), "user_id").Int()
	tokenType := gjson.Get(w.Body.String(), "token_type").String()
	accessToken := gjson.Get(w.Body.String(), "access_token").String()

	// Test results
	st.Expect(t, w.Code, 200)
	st.Expect(t, tokenType, "bearer")
	st.Expect(t, uint(userId), u1.Id)
	st.Expect(t, len(accessToken), 50)

	// Look in the sessions table for the correct access token
	sess, err := db.GetByAccessToken(accessToken)
	st.Expect(t, err, nil)
	st.Expect(t, sess.UserId, u1.Id)
	st.Expect(t, sess.ApplicationId, app.Id)
	st.Expect(t, sess.AccessToken, accessToken)

	// Let's get the user to make sure the accounts are added
	u, err := db.GetUserById(u1.Id)
	st.Expect(t, err, nil)
	st.Expect(t, len(u.Accounts), 1)
	st.Expect(t, u.Accounts[0].Name, account.Name)
}

//
// TestDoResetPassword02 - Validate password.
//
func TestDoResetPassword02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test users.
	u1 := test.GetRandomUser(33)
	db.Save(&u1)

	// Set reset password token
	token := "wootsblahdope"
	y := models.ForgotPassword{Token: token, UserId: u1.Id}
	db.Save(&y)

	// JSON we post in
	postStr := fmt.Sprintf(`{ "hash": "%s", "password": "123" }`, token)

	// Setup request
	req, _ := http.NewRequest("POST", "/reset-password", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.POST("/reset-password", c.DoResetPassword)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"Please enter a password at least 6 chars long."}`)
}

//
// TestDoResetPassword03 - Validate hash.
//
func TestDoResetPassword03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test users.
	u1 := test.GetRandomUser(33)
	db.Save(&u1)

	// Set reset password token
	token := "wootsblahdope"
	y := models.ForgotPassword{Token: token, UserId: u1.Id}
	db.Save(&y)

	// JSON we post in
	postStr := fmt.Sprintf(`{ "hash": "%s", "password": "helloworld123" }`, "wrongtoken")

	// Setup request
	req, _ := http.NewRequest("POST", "/reset-password", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.POST("/reset-password", c.DoResetPassword)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"Sorry, It seems your reset token has expired."}`)
}

//
// TestDoResetPassword04 - Validate hash. - empty token
//
func TestDoResetPassword04(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test users.
	u1 := test.GetRandomUser(33)
	db.Save(&u1)

	// Set reset password token
	token := "wootsblahdope"
	y := models.ForgotPassword{Token: token, UserId: u1.Id}
	db.Save(&y)

	// JSON we post in
	postStr := fmt.Sprintf(`{ "hash": "%s", "password": "helloworld123" }`, "")

	// Setup request
	req, _ := http.NewRequest("POST", "/reset-password", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.POST("/reset-password", c.DoResetPassword)
	r.ServeHTTP(w, req)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, w.Body.String(), `{"error":"Sorry, It seems your reset token has expired."}`)
}

/* End File */
