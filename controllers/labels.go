//
// Date: 2018-03-21
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-22
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"github.com/cloudmanic/skyclerk.com/library/request"
	"github.com/cloudmanic/skyclerk.com/library/response"
	"github.com/cloudmanic/skyclerk.com/models"
	"github.com/gin-gonic/gin"
)

//
// Return a list of labels
//
func (t *Controller) GetLabels(c *gin.Context) {

	// Place to store the results.
	var results = []models.Label{}

	// Get limits and pages
	page, _, _ := request.GetSetPagingParms(c)

	// Set the query parms
	params := models.QueryParam{
		Order:            c.DefaultQuery("order", "LabelsName"),
		Sort:             c.DefaultQuery("sort", "ASC"),
		Limit:            defaultLimit,
		Page:             page,
		AllowedOrderCols: []string{"LabelsId", "LabelsName"},
		Wheres: []models.KeyValue{
			{Key: "LabelsAccountId", ValueInt: c.MustGet("account").(int)},
		},
	}

	// Run the query
	err := t.db.Query(&results, params)

	// Get no filter count.
	noFilterCount, _ := t.db.QueryWithNoFilterCount(&results, params)

	// Get the meta data related to this query.
	meta := t.db.GetQueryMetaData(len(results), noFilterCount, params)

	// Put meta data in header.
	response.AddPagingInfoToHeaders(c, meta)

	// Return json based on if this was a good result or not.
	response.Results(c, results, err)
}

/* End File */
