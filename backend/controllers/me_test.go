//
// Date: 2019-04-14
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

	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
	"golang.org/x/crypto/bcrypt"

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
		c.Set("userId", int(user.Id))
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

//
// TestChangePassword01 - change user password
//
func TestChangePassword01(t *testing.T) {
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

	// Setup post string - F00bAr123 comes from the testing library
	postStr := fmt.Sprintf(`{ "current": "%s", "password": "%s", "confirm": "%s" }`, "F00bAr123", "foobar2", "foobar2")

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/me/change-password", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.POST("/api/v3/me/change-password", c.ChangePassword)
	r.ServeHTTP(w, req)

	// Get user.
	u := models.User{}
	db.Find(&u, 1)

	// Test results
	st.Expect(t, w.Code, 204)
	st.Expect(t, len(u.Md5Salt), 0)
	st.Expect(t, len(u.Md5Password), 0)
	st.Expect(t, len(u.Password) > 0, true)
}

//
// TestChangePassword02 - change user password but not MD5 based.
//
func TestChangePassword02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// New hash
	hash, _ := bcrypt.GenerateFromPassword([]byte("F00bAr123"), bcrypt.DefaultCost)

	// Setup test data
	user := test.GetRandomUser(33)
	user.Md5Salt = ""
	user.Md5Password = ""
	user.Password = string(hash)
	db.Save(&user)

	account1 := test.GetRandomAccount(33)
	account1.OwnerId = user.Id
	db.Save(&account1)
	db.Save(&models.AcctToUsers{AcctId: account1.Id, UserId: user.Id})

	// Setup post string - F00bAr123 comes from the testing library
	postStr := fmt.Sprintf(`{ "current": "%s", "password": "%s", "confirm": "%s" }`, "F00bAr123", "foobar2", "foobar2")

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/me/change-password", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.POST("/api/v3/me/change-password", c.ChangePassword)
	r.ServeHTTP(w, req)

	// Get user.
	u := models.User{}
	db.Find(&u, 1)

	// Test results
	st.Expect(t, w.Code, 204)
	st.Expect(t, len(u.Md5Salt), 0)
	st.Expect(t, len(u.Md5Password), 0)
	st.Expect(t, len(u.Password) > 0, true)
}

//
// TestChangePassword03 - Wrong password  - non-md5
//
func TestChangePassword03(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// New hash
	hash, _ := bcrypt.GenerateFromPassword([]byte("F00bAr123"), bcrypt.DefaultCost)

	// Setup test data
	user := test.GetRandomUser(33)
	user.Md5Salt = ""
	user.Md5Password = ""
	user.Password = string(hash)
	db.Save(&user)

	account1 := test.GetRandomAccount(33)
	account1.OwnerId = user.Id
	db.Save(&account1)
	db.Save(&models.AcctToUsers{AcctId: account1.Id, UserId: user.Id})

	// Setup post string - F00bAr123 comes from the testing library
	postStr := fmt.Sprintf(`{ "current": "%s", "password": "%s", "confirm": "%s" }`, "wrong", "foobar2", "foobar2")

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/me/change-password", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.POST("/api/v3/me/change-password", c.ChangePassword)
	r.ServeHTTP(w, req)

	// Get user.
	u := models.User{}
	db.Find(&u, 1)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, u.Password, string(hash))
	st.Expect(t, w.Body.String(), `{"error":"Your current password is not correct."}`)
}

//
// TestChangePassword04 - Wrong password  - md5 version
//
func TestChangePassword04(t *testing.T) {
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

	// Setup post string - F00bAr123 comes from the testing library
	postStr := fmt.Sprintf(`{ "current": "%s", "password": "%s", "confirm": "%s" }`, "wrong", "foobar2", "foobar2")

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/me/change-password", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.POST("/api/v3/me/change-password", c.ChangePassword)
	r.ServeHTTP(w, req)

	// Get user.
	u := models.User{}
	db.Find(&u, 1)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, u.Md5Salt, user.Md5Salt)
	st.Expect(t, u.Md5Password, user.Md5Password)
	st.Expect(t, u.Password, "")
	st.Expect(t, w.Body.String(), `{"error":"Your current password is not correct."}`)
}

//
// TestChangePassword05 - non-matching passwords
//
func TestChangePassword05(t *testing.T) {
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

	// Setup post string - F00bAr123 comes from the testing library
	postStr := fmt.Sprintf(`{ "current": "%s", "password": "%s", "confirm": "%s" }`, "F00bAr123", "match 1", "match 2")

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/me/change-password", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.POST("/api/v3/me/change-password", c.ChangePassword)
	r.ServeHTTP(w, req)

	// Get user.
	u := models.User{}
	db.Find(&u, 1)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, u.Md5Salt, user.Md5Salt)
	st.Expect(t, u.Md5Password, user.Md5Password)
	st.Expect(t, u.Password, "")
	st.Expect(t, w.Body.String(), `{"error":"Passwords do not match."}`)
}

/* End File */
