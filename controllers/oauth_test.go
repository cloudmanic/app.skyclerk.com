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

	"github.com/cloudmanic/skyclerk.com/library/test"
	"github.com/cloudmanic/skyclerk.com/models"
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
// TestDoOauthToken01 test to make sure auth works.
//
func TestDoOauthToken01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("testing_db")
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

/* End File */
