//
// Date: 11/3/2019
// Author(s): Spicer Matthews (spicer@options.cafe)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package admin

import (
	"github.com/gin-gonic/gin"
)

//
// DoRoutes sets the admin routes
//
func (t *Controller) DoRoutes(r *gin.Engine) {

	// ------------- Admin API --------------- //

	adminAPI := r.Group("/api/admin")

	adminAPI.Use(t.AuthMiddleware())
	{
		// Ping
		adminAPI.GET("/ping", t.PingFromServer)

		// Accounts
		adminAPI.GET("/accounts", t.GetAccounts)

		// Contacts
		adminAPI.GET("/contacts", t.GetContacts)

		// Categories
		adminAPI.GET("/categories", t.GetCategories)

		// Snapclerk
		adminAPI.GET("/snapclerk", t.GetSnapClerks)
		adminAPI.POST("/snapclerk/reject/:id", t.RejectSnapClerk)
		adminAPI.POST("/snapclerk/convert/:id", t.ConvertSnapClerk)

		// Users
		// adminApi.GET("/users", t.GetUsers)
		adminAPI.GET("/users/:id", t.GetUser)
		// adminApi.POST("/users/login-as-user", t.LoginAsUser)
	}

}

/* End File */
