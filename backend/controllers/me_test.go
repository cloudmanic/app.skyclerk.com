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

	account1 := test.GetRandomAccount(33)
	account1.OwnerId = user.Id
	account1.Name = "Matthews Etc."
	db.Save(&account1)
	user.Accounts = append(user.Accounts, account1)

	account2 := test.GetRandomAccount(34)
	account2.OwnerId = user.Id
	account2.Name = "Cloudmanic Labs, LLC"
	db.Save(&account2)
	user.Accounts = append(user.Accounts, account2)

	account3 := test.GetRandomAccount(105)
	account3.OwnerId = user.Id
	account3.Name = "124 West Main Street"
	db.Save(&account3)
	user.Accounts = append(user.Accounts, account3)

	// Save user.
	db.Save(&user)

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
	st.Expect(t, result.Accounts[0].Name, "124 West Main Street")
	st.Expect(t, result.Accounts[1].Name, "Cloudmanic Labs, LLC")
	st.Expect(t, result.Accounts[2].Name, "Matthews Etc.")
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
	db.Save(&models.AcctToUsers{AccountId: account1.Id, UserId: user.Id})

	// Setup post string - F00bAr123 comes from the testing library
	postStr := fmt.Sprintf(`{ "current": "%s", "password": "%s", "confirm": "%s" }`, "F00bAr123", "foobar2", "foobar2")

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/me/change-password", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.POST("/api/v3/:account/me/change-password", c.ChangePassword)
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
	db.Save(&models.AcctToUsers{AccountId: account1.Id, UserId: user.Id})

	// Setup post string - F00bAr123 comes from the testing library
	postStr := fmt.Sprintf(`{ "current": "%s", "password": "%s", "confirm": "%s" }`, "F00bAr123", "foobar2", "foobar2")

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/me/change-password", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.POST("/api/v3/:account/me/change-password", c.ChangePassword)
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
	db.Save(&models.AcctToUsers{AccountId: account1.Id, UserId: user.Id})

	// Setup post string - F00bAr123 comes from the testing library
	postStr := fmt.Sprintf(`{ "current": "%s", "password": "%s", "confirm": "%s" }`, "wrong", "foobar2", "foobar2")

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/me/change-password", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.POST("/api/v3/:account/me/change-password", c.ChangePassword)
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
	db.Save(&models.AcctToUsers{AccountId: account1.Id, UserId: user.Id})

	// Setup post string - F00bAr123 comes from the testing library
	postStr := fmt.Sprintf(`{ "current": "%s", "password": "%s", "confirm": "%s" }`, "wrong", "foobar2", "foobar2")

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/me/change-password", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.POST("/api/v3/:account/me/change-password", c.ChangePassword)
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
	db.Save(&models.AcctToUsers{AccountId: account1.Id, UserId: user.Id})

	// Setup post string - F00bAr123 comes from the testing library
	postStr := fmt.Sprintf(`{ "current": "%s", "password": "%s", "confirm": "%s" }`, "F00bAr123", "match 1", "match 2")

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/me/change-password", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.POST("/api/v3/:account/me/change-password", c.ChangePassword)
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

//
// TestChangePassword06 - validate password for strongness
//
func TestChangePassword06(t *testing.T) {
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
	db.Save(&models.AcctToUsers{AccountId: account1.Id, UserId: user.Id})

	// Setup post string - F00bAr123 comes from the testing library
	postStr := fmt.Sprintf(`{ "current": "%s", "password": "%s", "confirm": "%s" }`, "F00bAr123", "1", "1")

	// Setup request
	req, _ := http.NewRequest("POST", "/api/v3/33/me/change-password", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.POST("/api/v3/:account/me/change-password", c.ChangePassword)
	r.ServeHTTP(w, req)

	// Get user.
	u := models.User{}
	db.Find(&u, 1)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, u.Md5Salt, user.Md5Salt)
	st.Expect(t, u.Md5Password, user.Md5Password)
	st.Expect(t, u.Password, "")
	st.Expect(t, w.Body.String(), `{"error":"The password filed must be at least 6 characters long."}`)
}

//
// TestUpdateMe01 - update user profile
//
func TestUpdateMe01(t *testing.T) {
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
	db.Save(&models.AcctToUsers{AccountId: account1.Id, UserId: user.Id})

	// Update user profile
	postStr := fmt.Sprintf(`{ "first_name": "%s", "last_name": "%s", "email": "%s" }`, "Jane", "Wells", "jane@woots.com")

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/me", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.PUT("/api/v3/:account/me", c.UpdateMe)
	r.ServeHTTP(w, req)

	// Get user.
	u := models.User{}
	db.Find(&u, 1)

	// Test results
	st.Expect(t, w.Code, 204)
	st.Expect(t, u.FirstName, "Jane")
	st.Expect(t, u.LastName, "Wells")
	st.Expect(t, u.Email, "jane@woots.com")
}

//
// TestUpdateMe02 - update user profile - conflict in email
//
func TestUpdateMe02(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Setup test data
	user := test.GetRandomUser(33)
	db.Save(&user)

	// Setup test data
	user2 := test.GetRandomUser(33)
	db.Save(&user2)

	account1 := test.GetRandomAccount(33)
	account1.OwnerId = user.Id
	db.Save(&account1)
	db.Save(&models.AcctToUsers{AccountId: account1.Id, UserId: user.Id})

	// Update user profile
	postStr := fmt.Sprintf(`{ "first_name": "%s", "last_name": "%s", "email": "%s" }`, "Jane", "Wells", user2.Email)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/me", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.PUT("/api/v3/:account/me", c.UpdateMe)
	r.ServeHTTP(w, req)

	// Get user.
	u := models.User{}
	db.Find(&u, 1)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, u.FirstName, user.FirstName)
	st.Expect(t, u.LastName, user.LastName)
	st.Expect(t, u.Email, user.Email)
	st.Expect(t, w.Body.String(), `{"error":"Email already in use."}`)
}

//
// TestUpdateMe03 - make sure there are no issues with not updateing the email address
//
func TestUpdateMe03(t *testing.T) {
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
	db.Save(&models.AcctToUsers{AccountId: account1.Id, UserId: user.Id})

	// Update user profile
	postStr := fmt.Sprintf(`{ "first_name": "%s", "last_name": "%s", "email": "%s" }`, "Jane", "Wells", user.Email)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/me", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.PUT("/api/v3/:account/me", c.UpdateMe)
	r.ServeHTTP(w, req)

	// Get user.
	u := models.User{}
	db.Find(&u, 1)

	// Test results
	st.Expect(t, w.Code, 204)
	st.Expect(t, u.FirstName, "Jane")
	st.Expect(t, u.LastName, "Wells")
	st.Expect(t, u.Email, user.Email)
}

//
// TestUpdateMe04 - update user profile - Make sure we have a first name
//
func TestUpdateMe04(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Setup test data
	user := test.GetRandomUser(33)
	db.Save(&user)

	// Setup test data
	user2 := test.GetRandomUser(33)
	db.Save(&user2)

	account1 := test.GetRandomAccount(33)
	account1.OwnerId = user.Id
	db.Save(&account1)
	db.Save(&models.AcctToUsers{AccountId: account1.Id, UserId: user.Id})

	// Update user profile
	postStr := fmt.Sprintf(`{ "last_name": "%s", "email": "%s" }`, "Wells", user2.Email)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/me", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.PUT("/api/v3/:account/me", c.UpdateMe)
	r.ServeHTTP(w, req)

	// Get user.
	u := models.User{}
	db.Find(&u, 1)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, u.FirstName, user.FirstName)
	st.Expect(t, u.LastName, user.LastName)
	st.Expect(t, u.Email, user.Email)
	st.Expect(t, w.Body.String(), `{"error":"First name field is required."}`)
}

//
// TestUpdateMe05 - update user profile - Make sure we have a last name
//
func TestUpdateMe05(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Setup test data
	user := test.GetRandomUser(33)
	db.Save(&user)

	// Setup test data
	user2 := test.GetRandomUser(33)
	db.Save(&user2)

	account1 := test.GetRandomAccount(33)
	account1.OwnerId = user.Id
	db.Save(&account1)
	db.Save(&models.AcctToUsers{AccountId: account1.Id, UserId: user.Id})

	// Update user profile
	postStr := fmt.Sprintf(`{ "first_name": "%s", "email": "%s" }`, "Wells", user2.Email)

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/me", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.PUT("/api/v3/:account/me", c.UpdateMe)
	r.ServeHTTP(w, req)

	// Get user.
	u := models.User{}
	db.Find(&u, 1)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, u.FirstName, user.FirstName)
	st.Expect(t, u.LastName, user.LastName)
	st.Expect(t, u.Email, user.Email)
	st.Expect(t, w.Body.String(), `{"error":"Last name field is required."}`)
}

//
// TestUpdateMe06 - update user profile - Make sure we have an email
//
func TestUpdateMe06(t *testing.T) {
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Setup test data
	user := test.GetRandomUser(33)
	db.Save(&user)

	// Setup test data
	user2 := test.GetRandomUser(33)
	db.Save(&user2)

	account1 := test.GetRandomAccount(33)
	account1.OwnerId = user.Id
	db.Save(&account1)
	db.Save(&models.AcctToUsers{AccountId: account1.Id, UserId: user.Id})

	// Update user profile
	postStr := fmt.Sprintf(`{ "first_name": "%s", "last_name": "%s" }`, "Wells", "Wells")

	// Setup request
	req, _ := http.NewRequest("PUT", "/api/v3/33/me", bytes.NewBuffer([]byte(postStr)))

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", int(user.Id))
	})
	r.PUT("/api/v3/:account/me", c.UpdateMe)
	r.ServeHTTP(w, req)

	// Get user.
	u := models.User{}
	db.Find(&u, 1)

	// Test results
	st.Expect(t, w.Code, 400)
	st.Expect(t, u.FirstName, user.FirstName)
	st.Expect(t, u.LastName, user.LastName)
	st.Expect(t, u.Email, user.Email)
	st.Expect(t, w.Body.String(), `{"error":"Email field is required."}`)
}

/* End File */
