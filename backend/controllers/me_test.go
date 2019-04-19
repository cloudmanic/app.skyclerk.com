//
// Date: 2019-04-14
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nbio/st"

	"app.skyclerk.com/backend/library/test"
	"app.skyclerk.com/backend/models"
)

//
// TestGetMe01 - test getting me
//
func TestGetMe01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Setup test data
	user := test.GetRandomUser(33)
	db.Save(&user)

	account1 := test.GetRandomAccount(33)
	account1.OwnerId = user.Id
	db.Save(&account1)
	db.Save(&models.AcctToUsers{AcctId: account1.Id, UserId: user.Id})

	account2 := test.GetRandomAccount(34)
	account2.OwnerId = user.Id
	db.Save(&account2)
	db.Save(&models.AcctToUsers{AcctId: account2.Id, UserId: user.Id})

	account3 := test.GetRandomAccount(105)
	account3.OwnerId = user.Id
	db.Save(&account3)
	db.Save(&models.AcctToUsers{AcctId: account3.Id, UserId: user.Id})

	// Setup request
	req, _ := http.NewRequest("GET", "/oauth/me", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userId", user.Id)
	})
	r.GET("/oauth/me", c.GetMe)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	result := models.User{}
	err := json.Unmarshal([]byte(w.Body.String()), &result)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, result.Id, uint(1))
	st.Expect(t, result.FirstName, user.FirstName)
	st.Expect(t, result.LastName, user.LastName)
	st.Expect(t, result.Email, user.Email)
	st.Expect(t, len(result.Accounts), 3)
}

/* End File */
