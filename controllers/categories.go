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
// Return a list of categories. We limit to 500 mainly so we do not overload the
// system, but enough so the front-end does not have to page
//
func (t *Controller) GetCategories(c *gin.Context) {

	// Place to store the results.
	var results = []models.Category{}

	// Get limits and pages
	page, _, _ := request.GetSetPagingParms(c)

	// Set the query parms
	params := models.QueryParam{
		Order:            c.DefaultQuery("order", "CategoriesName"),
		Sort:             c.DefaultQuery("sort", "ASC"),
		Limit:            500,
		Page:             page,
		AllowedOrderCols: []string{"CategoriesId", "CategoriesName"},
		Wheres: []models.KeyValue{
			{Key: "CategoriesAccountId", ValueInt: c.MustGet("account").(int)},
		},
	}

	// Did we pass in a type of income?
	if c.Query("type") == "income" {
		params.Wheres = append(params.Wheres, models.KeyValue{Key: "CategoriesType", Value: "2"})
	}

	// Did we pass in a type of expense?
	if c.Query("type") == "expense" {
		params.Wheres = append(params.Wheres, models.KeyValue{Key: "CategoriesType", Value: "1"})
	}

	// Run the query
	meta, err := t.db.QueryMeta(&results, params)

	// Make the API more clear TODO: get rid of the numbering in the db in the future once we kill PHP
	for key, row := range results {
		if row.Type == "1" {
			results[key].Type = "expense"
		} else {
			results[key].Type = "income"
		}
	}

	// Return json based on if this was a good result or not.
	response.ResultsMeta(c, results, err, meta)
}

/* End File */
