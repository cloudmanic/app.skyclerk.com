//
// Date: 2019-09-13
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"app.skyclerk.com/backend/library/test"
	"app.skyclerk.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
	"golang.org/x/crypto/bcrypt"
)

// TestDoRegister01 Test registring a new user.
func TestDoRegister01(t *testing.T) {
	// Supply an isolated successful challenge for existing registration scenarios.
	mockRegisterTurnstile(t)
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
	postStr := fmt.Sprintf(`{ "turnstile_token": "test-token", "first": "%s", "last": "%s", "email": "%s", "password": "%s", "client_id": "%s" }`, "Jane", "Wells", "jane@wells.com", "foobar123", app.ClientId)

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
	db.Where("account_id = ? AND user_id = ?", 1, 1).First(&u)

	// Check the database that proper entries where created
	s := models.Session{}
	db.Where("user_id = ? AND application_id = ?", 1, 1).First(&s)

	// Check the database that proper entries where created
	m := models.User{}
	db.Where("id = ?", 1).First(&m)

	// Check the database that proper entries where created
	a := models.Account{}
	db.Where("owner_id = ?", 1).First(&a)

	// Check the database that proper entries where created
	b := models.Billing{}
	db.Where("id = ?", 1).First(&b)

	// Check the categories in the DB.
	cats := []models.Category{}
	db.Where("CategoriesAccountId = ?", 1).Find(&cats)

	// Test results
	st.Expect(t, w.Code, 200)
	st.Expect(t, res.UserId, uint(1))
	st.Expect(t, res.AccountId, uint(1))
	st.Expect(t, res.AccessToken, s.AccessToken)
	st.Expect(t, s.Id, uint(1))
	st.Expect(t, m.Id, uint(1))
	st.Expect(t, m.FirstName, "Jane")
	st.Expect(t, m.LastName, "Wells")
	st.Expect(t, m.Email, "jane@wells.com")
	st.Expect(t, a.Name, "Jane's Skyclerk")
	st.Expect(t, b.Id, uint(1))
	st.Expect(t, a.BillingId, uint(1))
	st.Expect(t, len(cats), 23)

	// Test password.
	err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte("foobar123"))
	st.Expect(t, err, nil)
}

// TestDoRegister02 - Error 01 (bad email)
func TestDoRegister02(t *testing.T) {
	// Supply an isolated successful challenge for existing registration scenarios.
	mockRegisterTurnstile(t)
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
	postStr := fmt.Sprintf(`{ "turnstile_token": "test-token", "first": "%s", "last": "%s", "email": "%s", "password": "%s", "client_id": "%s" }`, "Jane", "Wells", "jane", "foobar123", app.ClientId)

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

// TestDoRegister03 - Error 02 (no first)
func TestDoRegister03(t *testing.T) {
	// Supply an isolated successful challenge for existing registration scenarios.
	mockRegisterTurnstile(t)
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
	postStr := fmt.Sprintf(`{ "turnstile_token": "test-token", "first": "%s", "last": "%s", "email": "%s", "password": "%s", "client_id": "%s" }`, "", "Wells", "jane@wells.com", "foobar123", app.ClientId)

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

// TestDoRegister04 - Error 03 (bad password)
func TestDoRegister04(t *testing.T) {
	// Supply an isolated successful challenge for existing registration scenarios.
	mockRegisterTurnstile(t)
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
	postStr := fmt.Sprintf(`{ "turnstile_token": "test-token", "first": "%s", "last": "%s", "email": "%s", "password": "%s", "client_id": "%s" }`, "Jane", "Wells", "jane@wells.com", "ff", app.ClientId)

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

// TestDoRegister05 - Error 03 (no last)
func TestDoRegister05(t *testing.T) {
	// Supply an isolated successful challenge for existing registration scenarios.
	mockRegisterTurnstile(t)
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
	postStr := fmt.Sprintf(`{ "turnstile_token": "test-token", "first": "%s", "last": "%s", "email": "%s", "password": "%s", "client_id": "%s"  }`, "Jane", "", "jane@wells.com", "foobar123", app.ClientId)

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

// TestDoRegister06 - Error 04 (bad client id)
func TestDoRegister06(t *testing.T) {
	// Supply an isolated successful challenge for existing registration scenarios.
	mockRegisterTurnstile(t)
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
	postStr := fmt.Sprintf(`{ "turnstile_token": "test-token", "first": "%s", "last": "%s", "email": "%s", "password": "%s", "client_id": "%s"  }`, "Jane", "", "jane@wells.com", "foobar123", "bad")

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
	st.Expect(t, w.Body.String(), `{"error":"Something went wrong while logging into your account. Please try again or contact help@skyclerk.com. Sorry for the trouble."}`)
}

// TestDoRegister07 - Error 05 (missing client id)
func TestDoRegister07(t *testing.T) {
	// Supply an isolated successful challenge for existing registration scenarios.
	mockRegisterTurnstile(t)
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
	postStr := fmt.Sprintf(`{ "turnstile_token": "test-token", "first": "%s", "last": "%s", "email": "%s", "password": "%s" }`, "Jane", "", "jane@wells.com", "foobar123")

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
	st.Expect(t, w.Body.String(), `{"error":"Something went wrong while logging into your account. Please try again or contact help@skyclerk.com. Sorry for the trouble."}`)
}

// TestDoRegister08 - Error user already in the system.
func TestDoRegister08(t *testing.T) {
	// Supply an isolated successful challenge for existing registration scenarios.
	mockRegisterTurnstile(t)
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
	postStr := fmt.Sprintf(`{ "turnstile_token": "test-token", "first": "%s", "last": "%s", "email": "%s", "password": "%s", "client_id": "%s" }`, "Jane", "Wells", "jane@wells.com", "foobar123", app.ClientId)

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

// TestDoRegister09 Test registring a new user. With company name.
func TestDoRegister09(t *testing.T) {
	// Supply an isolated successful challenge for existing registration scenarios.
	mockRegisterTurnstile(t)
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
	postStr := fmt.Sprintf(`{ "turnstile_token": "test-token", "first": "%s", "last": "%s", "email": "%s", "password": "%s", "client_id": "%s", "company": "%s" }`, "Jane", "Wells", "jane@wells.com", "foobar123", app.ClientId, "ABC Inc.")

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
	db.Where("account_id = ? AND user_id = ?", 1, 1).First(&u)

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
	st.Expect(t, s.Id, uint(1))
	st.Expect(t, m.Id, uint(1))
	st.Expect(t, m.FirstName, "Jane")
	st.Expect(t, m.LastName, "Wells")
	st.Expect(t, m.Email, "jane@wells.com")
	st.Expect(t, a.Name, "ABC Inc.")

	// Test password.
	err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte("foobar123"))
	st.Expect(t, err, nil)
}

// TestDoRegister10 Test registring a new user. With token.
func TestDoRegister10(t *testing.T) {
	// Supply an isolated successful challenge for existing registration scenarios.
	mockRegisterTurnstile(t)
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create applicaiton.
	app := test.GetRandomApplication()
	app.GrantType = "password"
	db.Save(&app)

	// Create Accounts.
	acct1 := test.GetRandomAccount(33)
	acct2 := test.GetRandomAccount(27)
	db.Save(&acct1)
	db.Save(&acct2)

	// Expire
	now := time.Now()
	tExpire := now.Add(time.Hour * 24 * time.Duration(7))

	// Create an invite token
	invite := models.Invite{
		AccountId: acct2.Id,
		Email:     "katie@miller.com",
		FirstName: "Katie",
		LastName:  "Miller",
		Message:   "Woots, this is a message",
		Token:     "abc123",
		ExpiresAt: tExpire,
	}
	db.Save(&invite)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Get JSON
	postStr := fmt.Sprintf(`{ "turnstile_token": "test-token", "first": "%s", "last": "%s", "email": "%s", "password": "%s", "client_id": "%s", "company": "%s", "token": "%s" }`, "Jane", "Wells", "jane@wells.com", "foobar123", app.ClientId, "ABC Inc.", invite.Token)

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
	db.Where("account_id = ? AND user_id = ?", acct2.Id, 1).First(&u)

	// Check the database that proper entries where created
	s := models.Session{}
	db.Where("user_id = ? AND application_id = ?", 1, 1).First(&s)

	// Check the database that proper entries where created
	m := models.User{}
	db.Where("id = ?", 1).First(&m)

	// Check the database that proper entries where created
	a := models.Account{}
	db.Where("owner_id = ?", 1).First(&a)

	// Check the database that proper entries where created
	i := models.Invite{}
	db.Where("id = ?", 1).First(&i)

	// Test results
	st.Expect(t, w.Code, 200)
	st.Expect(t, res.UserId, uint(1))
	st.Expect(t, res.AccountId, acct2.Id)
	st.Expect(t, res.AccessToken, s.AccessToken)
	st.Expect(t, s.Id, uint(1))
	st.Expect(t, m.Id, uint(1))
	st.Expect(t, m.FirstName, "Jane")
	st.Expect(t, m.LastName, "Wells")
	st.Expect(t, m.Email, "jane@wells.com")
	st.Expect(t, a.Name, acct2.Name)
	st.Expect(t, i.Id, uint(0))

	// Test password.
	err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte("foobar123"))
	st.Expect(t, err, nil)
}

// TestDoRegister11 Test registring a new user. With bad token.
func TestDoRegister11(t *testing.T) {
	// Supply an isolated successful challenge for existing registration scenarios.
	mockRegisterTurnstile(t)
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create applicaiton.
	app := test.GetRandomApplication()
	app.GrantType = "password"
	db.Save(&app)

	// Create Accounts.
	acct1 := test.GetRandomAccount(33)
	acct2 := test.GetRandomAccount(27)
	db.Save(&acct1)
	db.Save(&acct2)

	// Expire
	now := time.Now()
	tExpire := now.Add(time.Hour * 24 * time.Duration(7))

	// Create an invite token
	invite := models.Invite{
		AccountId: acct2.Id,
		Email:     "katie@miller.com",
		FirstName: "Katie",
		LastName:  "Miller",
		Message:   "Woots, this is a message",
		Token:     "abc123",
		ExpiresAt: tExpire,
	}
	db.Save(&invite)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Get JSON
	postStr := fmt.Sprintf(`{ "turnstile_token": "test-token", "first": "%s", "last": "%s", "email": "%s", "password": "%s", "client_id": "%s", "company": "%s", "token": "%s" }`, "Jane", "Wells", "jane@wells.com", "foobar123", app.ClientId, "ABC Inc.", "lllllBAD")

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
	st.Expect(t, w.Body.String(), `{"error":"Your invite token is not found."}`)
}

// TestDoRegister12 Test registring a new user. With expired token.
func TestDoRegister12(t *testing.T) {
	// Supply an isolated successful challenge for existing registration scenarios.
	mockRegisterTurnstile(t)
	// Start the db connection.
	db, dbName, _ := models.NewTestDB("")
	defer models.TestingTearDown(db, dbName)

	// Create applicaiton.
	app := test.GetRandomApplication()
	app.GrantType = "password"
	db.Save(&app)

	// Create Accounts.
	acct1 := test.GetRandomAccount(33)
	acct2 := test.GetRandomAccount(27)
	db.Save(&acct1)
	db.Save(&acct2)

	// Expire (1 hour back)
	now := time.Now()
	tExpire := now.Add(time.Duration(-60) * time.Minute)

	// Create an invite token
	invite := models.Invite{
		AccountId: acct2.Id,
		Email:     "katie@miller.com",
		FirstName: "Katie",
		LastName:  "Miller",
		Message:   "Woots, this is a message",
		Token:     "abc123",
		ExpiresAt: tExpire,
	}
	db.Save(&invite)

	// Create controller
	c := &Controller{}
	c.SetDB(db)

	// Get JSON
	postStr := fmt.Sprintf(`{ "turnstile_token": "test-token", "first": "%s", "last": "%s", "email": "%s", "password": "%s", "client_id": "%s", "company": "%s", "token": "%s" }`, "Jane", "Wells", "jane@wells.com", "foobar123", app.ClientId, "ABC Inc.", invite.Token)

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
	st.Expect(t, w.Body.String(), `{"error":"Your invite token is not found."}`)
}

// mockRegisterTurnstile supplies fake credentials and a successful Cloudflare
// response so registration tests exercise verification without external services.
func mockRegisterTurnstile(t *testing.T) {
	// Isolate configuration and mock only Cloudflare's client, leaving other HTTP clients alone.
	t.Setenv("TURNSTILE_SECRET_KEY", "test-secret")
	t.Setenv("SITE_DOMAIN", "app.skyclerk.com")
	mockTurnstileResponse(t, 200, `{"success":true,"action":"register","hostname":"app.skyclerk.com"}`)

}

// turnstileTestTransport adapts a function into an isolated HTTP transport.
type turnstileTestTransport func(*http.Request) (*http.Response, error)

// RoundTrip invokes the fake verifier without changing the global HTTP transport.
func (transport turnstileTestTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return transport(request)
}

// mockTurnstileResponse validates outgoing form fields and supplies a deterministic
// response or network error. The previous client is restored after each test.
func mockTurnstileResponse(t *testing.T, status int, body string) *int {
	previous := registerTurnstileClient
	calls := 0
	// Restore the production client even when an assertion fails.
	t.Cleanup(func() { registerTurnstileClient = previous })
	// Keep all verification traffic in-process and assert the endpoint and payload.
	registerTurnstileClient = &http.Client{Transport: turnstileTestTransport(func(request *http.Request) (*http.Response, error) {
		calls++
		st.Expect(t, request.URL.String(), "https://challenges.cloudflare.com/turnstile/v0/siteverify")
		st.Expect(t, request.Method, http.MethodPost)
		st.Expect(t, request.Header.Get("Content-Type"), "application/x-www-form-urlencoded")
		// Decode the POST fields to ensure neither secret nor token is omitted.
		if err := request.ParseForm(); err != nil {
			t.Fatal(err)
		}
		st.Expect(t, request.PostForm.Get("secret"), "test-secret")
		st.Expect(t, request.PostForm.Get("response"), "test-token")
		if status == -1 {
			return nil, errors.New("verification unavailable")
		}
		if status == 0 {
			t.Fatal("unexpected verification request")
		}
		return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}, nil
	})}
	return &calls
}

// TestDoRegisterTurnstile checks that unverified signups cannot create database
// records, including when the caller supplies a valid invitation.
func TestDoRegisterTurnstile(t *testing.T) {
	cases := []struct {
		name, token, body string
		status            int
	}{
		{"missing token", "", "", 0},
		{"blank token", "   ", "", 0},
		{"oversized token", strings.Repeat("a", 2049), "", 0},
		{"missing secret", "test-token", "", 0},
		{"invalid token", "test-token", `{"success":false}`, 200},
		{"expired or replayed", "test-token", `{"success":false,"error-codes":["timeout-or-duplicate"]}`, 200},
		{"wrong action", "test-token", `{"success":true,"action":"login","hostname":"app.skyclerk.com"}`, 200},
		{"wrong hostname", "test-token", `{"success":true,"action":"register","hostname":"other.example"}`, 200},
		{"missing hostname", "test-token", `{"success":true,"action":"register","hostname":"app.skyclerk.com"}`, 200},
		{"malformed JSON", "test-token", `{`, 200},
		{"service unavailable", "test-token", `{"success":true,"action":"register","hostname":"app.skyclerk.com"}`, 503},
		{"network failure", "test-token", "", -1},
	}
	for _, tc := range cases {
		// Exercise failures through the real controller with a fresh database.
		t.Run(tc.name, func(t *testing.T) {
			// Missing configuration must fail closed even in local mode.
			t.Setenv("APP_ENV", "local")
			t.Setenv("TURNSTILE_SECRET_KEY", "test-secret")
			t.Setenv("SITE_DOMAIN", "app.skyclerk.com")
			if tc.name == "missing secret" {
				t.Setenv("TURNSTILE_SECRET_KEY", "")
			}
			if tc.name == "missing hostname" {
				t.Setenv("SITE_DOMAIN", "")
			}
			// Mock only verification requests, avoiding global HTTP transport changes.
			calls := mockTurnstileResponse(t, tc.status, tc.body)
			// Seed a valid client, account and invitation to reach the security check.
			db, dbName, err := models.NewTestDB("")
			if err != nil {
				t.Fatal(err)
			}
			defer models.TestingTearDown(db, dbName)
			application := test.GetRandomApplication()
			application.GrantType = "password"
			db.Save(&application)
			account := test.GetRandomAccount(33)
			db.Save(&account)
			invite := models.Invite{AccountId: account.Id, Token: "valid-invite", ExpiresAt: time.Now().Add(time.Hour)}
			db.Save(&invite)
			// Submit an otherwise valid invited signup.
			body, _ := json.Marshal(map[string]string{"first": "Jane", "last": "Wells", "email": "jane@wells.com", "password": "foobar123", "client_id": application.ClientId, "token": invite.Token, "turnstile_token": tc.token})
			controller := &Controller{}
			controller.SetDB(db)
			router := gin.New()
			router.POST("/register", controller.DoRegister)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body)))
			// Require the security error without any user, session or billing side effects.
			st.Expect(t, writer.Code, http.StatusBadRequest)
			st.Expect(t, writer.Body.String(), `{"error":"Please complete the security check and try again."}`)
			for _, model := range []interface{}{&models.User{}, &models.Session{}, &models.Billing{}, &models.AcctToUsers{}} {
				var count int
				db.Model(model).Count(&count)
				st.Expect(t, count, 0)
			}
			var accounts, invites int
			db.Model(&models.Account{}).Count(&accounts)
			db.Model(&models.Invite{}).Count(&invites)
			st.Expect(t, accounts, 1)
			st.Expect(t, invites, 1)
			st.Expect(t, *calls > 0, tc.status != 0)
		})
	}
}

// TestRegisterConfig ensures the browser receives only the public key, with no
// caching, and registration is unavailable when either key is missing.
func TestRegisterConfig(t *testing.T) {
	for _, keys := range []struct {
		site, secret string
		status       int
	}{
		{"public-key", "private-key", 200}, {"", "private-key", 503}, {"public-key", "", 503},
	} {
		// Isolate each configuration from real runtime credentials.
		t.Run(fmt.Sprintf("%s-%d", keys.site, keys.status), func(t *testing.T) {
			t.Setenv("TURNSTILE_SITE_KEY", keys.site)
			t.Setenv("TURNSTILE_SECRET_KEY", keys.secret)
			controller := &Controller{}
			router := gin.New()
			router.GET("/registration-config", controller.RegisterConfig)
			writer := httptest.NewRecorder()
			router.ServeHTTP(writer, httptest.NewRequest(http.MethodGet, "/registration-config", nil))
			// Assert the secret never appears and key rotation is not hidden by caches.
			st.Expect(t, writer.Code, keys.status)
			st.Expect(t, writer.Header().Get("Cache-Control"), "no-store")
			st.Expect(t, strings.Contains(writer.Body.String(), "private-key"), false)
			if keys.status == 200 {
				st.Expect(t, writer.Body.String(), `{"site_key":"public-key"}`)
			}
		})
	}
}

/* End File */
