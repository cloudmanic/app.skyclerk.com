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

//
// DoRoutes - Do Routes
//
func (t *Controller) DoRoutes(r *gin.Engine) {

	// --------- API V1 sub-routes ----------- //

	apiV1 := r.Group("/api/v1")

	apiV1.Use(t.AuthMiddleware())
	{
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
		apiV1.POST("/:account/contacts", t.CreateContact)
	}

	// ------------ Non-Auth Routes ------ //

	// // Auth Routes
	//r.POST("/oauth/token", t.DoOauthToken)

	// -------- Static Files ------------ //

	r.Use(static.Serve("/", static.LocalFile("/frontend", true)))
	r.NoRoute(func(c *gin.Context) { c.File("/frontend/index.html") })
}

/* End File */
