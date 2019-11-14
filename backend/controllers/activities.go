//
// Date: 2019-06-26
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"
	"strconv"

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
		PreLoads:         []string{"User", "Ledger"},
		Wheres: []models.KeyValue{
			{Key: "account_id", Compare: "=", ValueInt: c.MustGet("accountId").(int)},
		},
	}

	// Did we pass in a leder_id so we filter by a ledger.
	if c.DefaultQuery("ledger_id", "") != "" {
		// Set id
		ledger_id, err := strconv.ParseInt(c.Query("ledger_id"), 10, 32)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err})
			return
		}

		params.Wheres = append(params.Wheres, models.KeyValue{
			Key:      "ledger_id",
			Compare:  "=",
			ValueInt: int(ledger_id),
		})
	}

	// Run the query
	meta, err := t.db.QueryMeta(&results, params)

	// Add in the message
	for key, row := range results {
		// Add in ledger contact
		if row.Ledger.ContactId > 0 {
			c, _ := t.db.GetContactByAccountAndId(row.AccountId, row.Ledger.ContactId)
			results[key].Ledger.Contact = c
		}

		// Set message
		results[key].SetMessage()
	}

	// Return json based on if this was a good result or not.
	if c.DefaultQuery("group", "") == "date" {
		r := make(map[string][]models.Activity)

		for _, row := range results {
			indx := row.CreatedAt.Format("2006-01-02")
			r[indx] = append(r[indx], row)
		}

		response.ResultsMeta(c, r, err, meta)
	} else {
		response.ResultsMeta(c, results, err, meta)
	}
}

/* End File */
