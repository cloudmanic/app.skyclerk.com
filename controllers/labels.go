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
// Return a list of labels. We limit to 500 mainly so we do not overload the
// system, but enough so the front-end does not have to page
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
		Limit:            500,
		Page:             page,
		AllowedOrderCols: []string{"LabelsId", "LabelsName"},
		Wheres: []models.KeyValue{
			{Key: "LabelsAccountId", ValueInt: c.MustGet("account").(int)},
		},
	}

	// Run the query
	meta, err := t.db.QueryMeta(&results, params)

	// Return json based on if this was a good result or not.
	response.ResultsMeta(c, results, err, meta)
}

/* End File */
