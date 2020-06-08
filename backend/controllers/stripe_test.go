//
// Date: 2020-05-07
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"app.skyclerk.com/backend/library/test"
	"app.skyclerk.com/backend/models"
	"github.com/gin-gonic/gin"
	"github.com/nbio/st"
)

//
// TestStripeAuthorizeURL01
//
func TestStripeAuthorizeURL01(t *testing.T) {
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

	// Setup request
	req, _ := http.NewRequest("GET", "/api/v3/33/stripe/authorize?redirect=https://yahoo.com", nil)

	// Setup writer.
	w := httptest.NewRecorder()
	gin.SetMode("release")
	gin.DisableConsoleColor()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("accountId", 33)
		c.Set("userId", uint(109))
	})
	r.GET("/api/v3/:account/stripe/authorize", c.StripeAuthorizeURL)
	r.ServeHTTP(w, req)

	// Test result TODO(spicer): verify redis is getting set with the correct value. Verify the rest of the url being returned.
	//redirectURL := fmt.Sprintf("%s/stripe/auth/callback", os.Getenv("APP_URL"))
	st.Expect(t, strings.Contains(w.Body.String(), "https://connect.stripe.com/oauth/authorize?"), true)
	st.Expect(t, strings.Contains(w.Body.String(), "client_id="+os.Getenv("STRIPE_CLIENT_ID")), true)
}

/* End File */
