//
// Date: 2018-03-21
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-12-29
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/cloudmanic/skyclerk.com/backend/library/request"
	"github.com/cloudmanic/skyclerk.com/backend/library/response"
	"github.com/cloudmanic/skyclerk.com/backend/models"
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
			{Key: "CategoriesAccountId", ValueInt: c.MustGet("accountId").(int)},
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

//
// Get Category by id
//
func (t *Controller) GetCategory(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	// Get category and make sure we have perms to it
	orgCat, err := t.db.GetCategoryByAccountAndId(uint(c.MustGet("accountId").(int)), uint(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found."})
		return
	}

	// Make the API more clear TODO: get rid of the numbering in the db in the future once we kill PHP
	if orgCat.Type == "1" {
		orgCat.Type = "expense"
	} else {
		orgCat.Type = "income"
	}

	// Return happy.
	response.Results(c, orgCat, nil)
}

//
// Create a category within the account.
//
func (t *Controller) CreateCategory(c *gin.Context) {

	// Setup Category obj
	o := models.Category{}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &o, "create") != nil {
		return
	}

	// Make sure the AccountId is correct.
	o.AccountId = uint(c.MustGet("accountId").(int))

	// Clean up some vars
	o.Type = strings.Trim(o.Type, " ")
	o.Name = strings.Trim(o.Name, " ")

	// Create category
	t.db.New().Create(&o)

	// Make the API more clear TODO: get rid of the numbering in the db in the future once we kill PHP
	if o.Type == "1" {
		o.Type = "expense"
	} else {
		o.Type = "income"
	}

	// Return happy.
	response.RespondCreated(c, o, nil)
}

//
// Update a category within the account.
//
func (t *Controller) UpdateCategory(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	// First we make sure this is an entry we have access to.
	orgCat, err := t.db.GetCategoryByAccountAndId(uint(c.MustGet("accountId").(int)), uint(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found."})
		return
	}

	// Setup Category obj
	o := models.Category{}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &o, "update") != nil {
		return
	}

	// We just allow updating of a few fields
	orgCat.Type = strings.Trim(o.Type, " ")
	orgCat.Name = strings.Trim(o.Name, " ")

	// Update category
	t.db.New().Save(&orgCat)

	// Make the API more clear TODO: get rid of the numbering in the db in the future once we kill PHP
	if orgCat.Type == "1" {
		orgCat.Type = "expense"
	} else {
		orgCat.Type = "income"
	}

	// Return happy.
	response.RespondUpdated(c, orgCat, nil)
}

//
// Delete a category within the account.
//
func (t *Controller) DeleteCategory(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	// First we make sure this is an entry we have access to.
	_, err2 := t.db.GetCategoryByAccountAndId(uint(c.MustGet("accountId").(int)), uint(id))

	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category not found."})
		return
	}

	// Delete category
	err = t.db.DeleteCategoryByAccountAndId(uint(c.MustGet("accountId").(int)), uint(id))

	if err != nil {
		response.RespondError(c, err)
		return
	}

	// Return happy.
	response.RespondDeleted(c, nil)
}

/* End File */
