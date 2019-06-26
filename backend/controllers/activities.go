//
// Date: 2019-06-26
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"app.skyclerk.com/backend/library/request"
	"app.skyclerk.com/backend/library/response"
	"app.skyclerk.com/backend/models"
	"github.com/gin-gonic/gin"
)

//
// GetActivities - Return a list of Activity. We limit to 100.
//
func (t *Controller) GetActivities(c *gin.Context) {
	// Place to store the results.
	var results = []models.Activity{}

	// Get limits and pages
	page, limit, _ := request.GetSetPagingParms(c)

	// Set the query parms
	params := models.QueryParam{
		Debug:            false,
		Order:            c.DefaultQuery("order", "id"),
		Sort:             c.DefaultQuery("sort", "DESC"),
		Limit:            limit,
		Page:             page,
		AllowedOrderCols: []string{},
		PreLoads:         []string{"User"},
		Wheres: []models.KeyValue{
			{Key: "user_id", Compare: "=", ValueInt: c.MustGet("userId").(int)},
			{Key: "account_id", Compare: "=", ValueInt: c.MustGet("accountId").(int)},
		},
	}

	// Run the query
	meta, err := t.db.QueryMeta(&results, params)

	// Add in the message
	for key := range results {
		results[key].SetMessage()
	}

	// Return json based on if this was a good result or not.
	response.ResultsMeta(c, results, err, meta)
}

/* End File */
