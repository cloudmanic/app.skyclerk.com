//
// Date: 2018-03-22
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-29
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"app.skyclerk.com/backend/library/test"
	"app.skyclerk.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// Test get a users 01
//
func TestGetUsers01(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Create test users.
	u1 := test.GetRandomUser(33)
	u2 := test.GetRandomUser(33)
	u3 := test.GetRandomUser(33)
	u4 := test.GetRandomUser(33)
	u5 := test.GetRandomUser(33)
	u6 := test.GetRandomUser(22)
	u7 := test.GetRandomUser(22)

	db.Save(&u1)
	db.Save(&u2)
	db.Save(&u3)
	db.Save(&u4)
	db.Save(&u5)
	db.Save(&u6)
	db.Save(&u7)

	db.Save(&models.AcctToUsers{AcctId: uint(33), UserId: u1.Id})
	db.Save(&models.AcctToUsers{AcctId: uint(33), UserId: u2.Id})
	db.Save(&models.AcctToUsers{AcctId: uint(33), UserId: u3.Id})
	db.Save(&models.AcctToUsers{AcctId: uint(33), UserId: u4.Id})
	db.Save(&models.AcctToUsers{AcctId: uint(33), UserId: u5.Id})
	db.Save(&models.AcctToUsers{AcctId: uint(22), UserId: u6.Id})
	db.Save(&models.AcctToUsers{AcctId: uint(22), UserId: u7.Id})

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/users", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/users", c.GetUsers)
	r.ServeHTTP(w, req)

	// Grab result and convert to strut
	results := []models.User{}
	err := json.Unmarshal([]byte(w.Body.String()), &results)

	// Test results
	st.Expect(t, err, nil)
	st.Expect(t, len(results), 5)
	st.Expect(t, results[0].Id, u1.Id)
	st.Expect(t, results[1].Id, u2.Id)
	st.Expect(t, results[2].Id, u3.Id)
	st.Expect(t, results[3].Id, u4.Id)
	st.Expect(t, results[4].Id, u5.Id)
	st.Expect(t, results[0].Email, u1.Email)
	st.Expect(t, results[1].Email, u2.Email)
	st.Expect(t, results[2].Email, u3.Email)
	st.Expect(t, results[3].Email, u4.Email)
	st.Expect(t, results[4].Email, u5.Email)
}

/* End File */
