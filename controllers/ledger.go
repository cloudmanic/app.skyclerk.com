//
// Date: 2019-01-13
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2019-01-13
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"github.com/gin-gonic/gin"

	"github.com/cloudmanic/skyclerk.com/library/request"
	"github.com/cloudmanic/skyclerk.com/library/response"
	"github.com/cloudmanic/skyclerk.com/models"
)

//
// GetLedgers - Return a list of ledgers. We limit to 50 mainly so we do not overload the
// system, but enough so the front-end does not have to page
//
func (t *Controller) GetLedgers(c *gin.Context) {
	// Place to store the results.
	var results = []models.Ledger{}

	// Get limits and pages
	page, _, _ := request.GetSetPagingParms(c)

	// Set the query parms
	params := models.QueryParam{
		Order:            c.DefaultQuery("order", "LedgerDate"),
		Sort:             c.DefaultQuery("sort", "DESC"),
		Limit:            50,
		Page:             page,
		Debug:            false,
		PreLoads:         []string{"Category", "Contact", "Labels"},
		AllowedOrderCols: []string{"LedgerId", "LedgerDate"},
		Wheres: []models.KeyValue{
			{Key: "LedgerAccountId", ValueInt: c.MustGet("accountId").(int)},
		},
	}

	// Run the query
	meta, err := t.db.QueryMeta(&results, params)

	// Return json based on if this was a good result or not.
	response.ResultsMeta(c, results, err, meta)
}

/* End File */
