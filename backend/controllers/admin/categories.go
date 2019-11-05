//
// Date: 11/3/2018
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package admin

import (
	"github.com/gin-gonic/gin"

	"app.skyclerk.com/backend/library/helpers"
	"app.skyclerk.com/backend/library/request"
	"app.skyclerk.com/backend/library/response"
	"app.skyclerk.com/backend/models"
)

//
// GetCategories Return a list of categories. We limit to 500 mainly so we do not overload the
// system, but enough so the front-end does not have to page
//
func (t *Controller) GetCategories(c *gin.Context) {
	// Get the account id
	accountID := helpers.StringToInt(c.DefaultQuery("account_id", "0"))

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
		Debug:            false,
		AllowedOrderCols: []string{"CategoriesId", "CategoriesName"},
		Wheres: []models.KeyValue{
			{Key: "CategoriesAccountId", Compare: "=", ValueInt: accountID},
		},
	}

	// Manage a search query - TODO(spicer): this is hacky but it will do for now.
	if len(c.DefaultQuery("search", "")) > 0 {
		// Query term
		search := c.DefaultQuery("search", "")

		// Array of category ids
		categoryWhereIn := []int{}

		// Get contacts
		con := []models.Category{}
		t.db.New().Where("(CategoriesName LIKE ?) AND CategoriesAccountId = ?", "%"+search+"%", accountID).Find(&con)

		// Build id array.
		for _, row := range con {
			categoryWhereIn = append(categoryWhereIn, int(row.Id))
		}

		// Update query.
		params.Wheres = append(params.Wheres, models.KeyValue{
			Key:          "CategoriesId",
			Compare:      "IN",
			ValueIntList: categoryWhereIn,
		})

		// If we have no ids we have no results
		if len(categoryWhereIn) == 0 {
			// Run the query
			meta, err := t.db.QueryMeta(&results, params)

			// Return json based on if this was a good result or not.
			response.ResultsMeta(c, []models.Category{}, err, meta)
			return
		}
	}

	// Did we pass in a type of income?
	if c.Query("type") == "income" {
		params.Wheres = append(params.Wheres, models.KeyValue{Key: "CategoriesType", Compare: "=", Value: "2"})
	}

	// Did we pass in a type of expense?
	if c.Query("type") == "expense" {
		params.Wheres = append(params.Wheres, models.KeyValue{Key: "CategoriesType", Compare: "=", Value: "1"})
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

	// Add in usage.
	usage := t.db.GetCategoryUsageByAccount(uint(accountID))

	// Add in count TODO(spicer): use MAPs so we do less loops
	for key, row := range results {
		for _, row2 := range usage {
			if row2.Name == row.Name {
				results[key].Count = row2.Count
			}
		}
	}

	// Return json based on if this was a good result or not.
	response.ResultsMeta(c, results, err, meta)
}

/* End File */
