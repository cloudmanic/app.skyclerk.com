//
// Date: 2020-05-07
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"app.skyclerk.com/backend/library/cache"
	"app.skyclerk.com/backend/library/helpers"
	"app.skyclerk.com/backend/library/slack"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v71/oauth"

	stripe "github.com/stripe/stripe-go/v71"
)

//
// StripeAuthorizeURL will return the authorize url for stripe.
//
func (t *Controller) StripeAuthorizeURL(c *gin.Context) {
	// Get account id
	accountID := fmt.Sprintf("%d", c.MustGet("accountId"))

	// Get account.
	_, err := t.db.GetAccountById(uint(c.MustGet("accountId").(int)))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found."})
		return
	}

	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"system": "No Stripe key setup."}})
		c.AbortWithStatus(401)
		return
	}

	// Setup State
	randomString := helpers.RandStr(20)

	// Store state in redis cache. Expires in 30 mins.
	cache.SetExpire("sc-stripe-callback-state-"+randomString, time.Duration(30)*time.Minute, string(accountID))

	// Setup redirect
	cache.SetExpire("sc-stripe-callback-redirect-"+randomString, time.Duration(30)*time.Minute, c.DefaultQuery("redirect", os.Getenv("SITE_URL")))

	// Redirect URL
	redirectURL := fmt.Sprintf("%s/stripe/auth/callback", os.Getenv("APP_URL"))

	// Build URL to redirect to.
	url := oauth.AuthorizeURL(&stripe.AuthorizeURLParams{
		ClientID:      stripe.String(os.Getenv("STRIPE_CLIENT_ID")),
		State:         stripe.String(randomString),
		Scope:         stripe.String("read_only"),
		RedirectURI:   stripe.String(redirectURL),
		ResponseType:  stripe.String("code"),
		StripeLanding: stripe.String("login"), // or "register"
		AlwaysPrompt:  stripe.Bool(true),      // Boolean to indicate that the user should always be asked to connect, even if they're already connected.
		Express:       stripe.Bool(false),
	})

	// Return happy JSON
	c.JSON(200, gin.H{"url": url})
}

//
// StripeAuthCallback is the call back stripe redirects to after an auth dance.
//
func (t *Controller) StripeAuthCallback(c *gin.Context) {
	// Set vars from url.
	code := c.Query("code")
	state := c.Query("state")
	cacheKey := "sc-stripe-callback-state-" + state
	redirectKey := "sc-stripe-callback-redirect-" + state

	// Get the account id from the state
	accountIDString := ""
	found, err := cache.Get(cacheKey, &accountIDString)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad state. Please contract help@skyclerk.com"})
		return
	}

	if !found {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad state. Please contract help@skyclerk.com"})
		return
	}

	// Convet to unt
	accountID := helpers.StringToUint(accountIDString)

	// Get account.
	account, err := t.db.GetAccountById(accountID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Account not found."})
		return
	}

	// Make sure we have a STRIPE_SECRET_KEY
	if len(os.Getenv("STRIPE_SECRET_KEY")) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No stripe key"})
		return
	}

	// Set the stripe key.
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Send the authorization code to Stripe's API.
	params := &stripe.OAuthTokenParams{
		GrantType: stripe.String("authorization_code"),
		Code:      &code,
	}

	// Get token.
	token, err := oauth.New(params)

	if err != nil {
		stripeErr := err.(*stripe.Error)
		if stripeErr.OAuthError == "invalid_grant" {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid authorization code: %s Please contract help@skyclerk.com", code)})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown error. Please contract help@skyclerk.com"})
		}
		return
	}

	ca, _ := t.db.GetConnectedAccountsByAccountIDAndConnection(account.Id, "Stripe")

	// Store in our connected table
	ca.AccountID = account.Id
	ca.Connection = "Stripe"
	ca.StripeUserID = token.StripeUserID
	ca.StripeAccessToken = token.AccessToken
	ca.StripeRefreshToken = token.RefreshToken
	ca.StripePublishableKey = token.StripePublishableKey
	ca.StripeScope = fmt.Sprintf("%s", token.Scope)

	t.db.New().Save(&ca)

	// Get the account owner
	user, _ := t.db.GetUserById(account.OwnerId)

	// Tell slack about this.
	go slack.Notify("#events", "Skyclerk - New Stripe Connected Account : "+user.Email)

	// Get the redirect URL
	redirectURL := ""
	found, _ = cache.Get(redirectKey, &redirectURL)

	// Clear the cache.
	cache.Delete(cacheKey)
	cache.Delete(redirectKey)

	if found {
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
	} else {
		c.Redirect(http.StatusTemporaryRedirect, os.Getenv("SITE_URL"))
	}
}

/* End File */
