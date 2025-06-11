//
// Date: 2018-03-21
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2019-01-13
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// DoRoutes - Do Routes
func (t *Controller) DoRoutes(r *gin.Engine) {

	// --------- API V1 sub-routes ----------- //

	apiV1 := r.Group("/api/v3")

	apiV1.Use(t.AuthMiddleware())
	{
		// Ping
		apiV1.GET("/:account/ping", t.PingFromClient)

		// Accounts
		apiV1.GET("/:account/account", t.GetAccount)
		apiV1.PUT("/:account/account", t.UpdateAccount)
		apiV1.POST("/:account/account/new", t.NewAccount)
		apiV1.POST("/:account/account/clear", t.ClearAccount)
		apiV1.POST("/:account/account/delete", t.DeleteAccount)
		apiV1.PUT("/:account/account/subscription", t.ChangeSubscription)
		apiV1.POST("/:account/account/stripe-token", t.NewStripeToken)
		apiV1.GET("/:account/account/billing", t.GetBilling)
		apiV1.GET("/:account/account/billing-history", t.GetBillingHistory)
		apiV1.POST("/:account/account/apple-in-app", t.UpdateAccountAppleInApp)

		// Me
		apiV1.PUT("/:account/me", t.UpdateMe)
		apiV1.POST("/:account/me/change-password", t.ChangePassword)

		// Users
		apiV1.GET("/:account/users", t.GetUsers)
		apiV1.GET("/:account/users/invite", t.GetInvitedUsers)
		apiV1.POST("/:account/users/invite", t.InviteUser)
		apiV1.DELETE("/:account/users/:id", t.DeleteUser)
		apiV1.DELETE("/:account/user-invite/:id", t.DeleteInvite)

		// Ledger
		apiV1.GET("/:account/ledger", t.GetLedgers)
		apiV1.GET("/:account/ledger/:id", t.GetLedger)
		apiV1.GET("/:account/ledger-summary", t.GetLedgerSummary)
		apiV1.GET("/:account/ledger-pl-summary", t.GetLedgerPlSummary)
		apiV1.POST("/:account/ledger", t.CreateLedger)
		apiV1.PUT("/:account/ledger/:id", t.UpdateLedger)
		apiV1.DELETE("/:account/ledger/:id", t.DeleteLedger)

		// Labels
		apiV1.GET("/:account/labels", t.GetLabels)
		apiV1.GET("/:account/labels/:id", t.GetLabel)
		apiV1.POST("/:account/labels", t.CreateLabel)
		apiV1.PUT("/:account/labels/:id", t.UpdateLabel)
		apiV1.DELETE("/:account/labels/:id", t.DeleteLabel)

		// Categories
		apiV1.GET("/:account/categories", t.GetCategories)
		apiV1.GET("/:account/categories/:id", t.GetCategory)
		apiV1.POST("/:account/categories", t.CreateCategory)
		apiV1.PUT("/:account/categories/:id", t.UpdateCategory)
		apiV1.DELETE("/:account/categories/:id", t.DeleteCategory)

		// Contacts
		apiV1.GET("/:account/contacts", t.GetContacts)
		apiV1.GET("/:account/contacts/:id", t.GetContact)
		apiV1.POST("/:account/contacts", t.CreateContact)
		apiV1.PUT("/:account/contacts/:id", t.UpdateContact)
		apiV1.DELETE("/:account/contacts/:id", t.DeleteContact)

		// Files
		apiV1.POST("/:account/files", t.CreateFile)

		// Snapclerk
		apiV1.GET("/:account/snapclerk", t.GetSnapClerk)
		apiV1.GET("/:account/snapclerk/usage", t.GetSnapClerkUsage)
		apiV1.POST("/:account/snapclerk", t.CreateSnapClerk)
		apiV1.POST("/:account/snapclerk/add-by-file-id", t.CreateSnapClerkByFileId)

		// Activities
		apiV1.GET("/:account/activities", t.GetActivities)

		// Reports
		apiV1.GET("/:account/reports/pnl", t.ReportsPnl)
		apiV1.GET("/:account/reports/pnl-label", t.ReportsPnlLabel)
		apiV1.GET("/:account/reports/pnl-category", t.ReportsPnlCategory)
		apiV1.GET("/:account/reports/income-by-contact", t.ReportsIncomeByContact)
		apiV1.GET("/:account/reports/expenses-by-contact", t.ReportsExpensesByContact)
		apiV1.GET("/:account/reports/pnl-current-year", t.ReportsCurrentPnl)

		// Stripe
		apiV1.GET("/:account/stripe/authorize", t.StripeAuthorizeURL)
	}

	// ------------ Non-API Routes ------ //

	// oauth Routes
	r.POST("/oauth/token", t.DoOauthToken)
	r.GET("/oauth/logout", t.DoLogOut)
	r.Group("/oauth/me").Use(t.AuthNoAccountMiddleware()).GET("", t.GetMe)

	// Other Auth Routes
	r.POST("/register", t.DoRegister)
	r.POST("/reset-password", t.DoResetPassword)
	r.POST("/forgot-password", t.DoForgotPassword)

	// Webhooks
	r.POST("/webhooks/fovea", t.DoFoveaWebhook)
	r.POST("/webhooks/postmark", t.DoPostmarkWebhook)

	// Support
	r.POST("/support/contact-us", t.ContactUs)

	// Stripe Auth Callback
	r.GET("/stripe/auth/callback", t.StripeAuthCallback)

	// -------- Static Files ------------ //

	r.Use(static.Serve("/", static.LocalFile("/app/frontend", true)))
	r.Use(static.Serve("/centcom", static.LocalFile("/app/centcom", true)))
	r.Use(static.Serve("/centcom/accounts", static.LocalFile("/app/centcom", true)))
	r.Use(static.Serve("/centcom/snapclerk", static.LocalFile("/app/centcom", true)))
	r.NoRoute(func(c *gin.Context) { c.File("/app/frontend/index.html") })
}

/* End File */
