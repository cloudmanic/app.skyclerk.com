//
// Date: 2019-01-13
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2019-01-13
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"app.skyclerk.com/backend/library/request"
	"app.skyclerk.com/backend/library/response"
	"app.skyclerk.com/backend/models"
	"app.skyclerk.com/backend/services"
)

//
// GetLedgers - Return a list of ledgers. We limit to 50 mainly so we do not overload the
// system, but enough so the front-end does not have to page
//
func (t *Controller) GetLedgers(c *gin.Context) {
	// Place to store the results.
	var results = []models.Ledger{}

	// Get account
	accountId := c.MustGet("accountId").(int)

	// Get limits and pages
	page, _, _ := request.GetSetPagingParms(c)

	// Set the query parms
	params := models.QueryParam{
		Order:            c.DefaultQuery("order", "LedgerDate"),
		Sort:             c.DefaultQuery("sort", "DESC"),
		Limit:            50,
		Page:             page,
		Debug:            false,
		PreLoads:         []string{"Category", "Contact", "Labels", "Files"},
		AllowedOrderCols: []string{"LedgerId", "LedgerDate"},
		Wheres: []models.KeyValue{
			{Key: "LedgerAccountId", Compare: "=", ValueInt: accountId},
		},
	}

	// Add type filter - income
	if c.DefaultQuery("type", "") == "income" {
		params.Wheres = append(params.Wheres, models.KeyValue{
			Key:        "LedgerAmount",
			Compare:    ">=",
			ValueFloat: 0.01, // because of the lib we can't use 0
		})
	}

	// Add type filter - expense
	if c.DefaultQuery("type", "") == "expense" {
		params.Wheres = append(params.Wheres, models.KeyValue{
			Key:        "LedgerAmount",
			Compare:    "<=",
			ValueFloat: -0.01, // because of the lib we can't use 0
		})
	}

	// Add type filter - category_id
	if len(c.DefaultQuery("category_id", "")) > 0 {
		// Convert cat id
		cat_id, err := strconv.Atoi(c.DefaultQuery("category_id", ""))

		if err != nil {
			services.Info(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error with category_id"})
			return
		}

		// Update query.
		params.Wheres = append(params.Wheres, models.KeyValue{
			Key:      "LedgerCategoryId",
			Compare:  "=",
			ValueInt: cat_id,
		})
	}

	// Get ledger ids from lables we want to filter from.
	if len(c.DefaultQuery("label_ids", "")) > 0 {
		// Get Ids from url
		ids := strings.Split(c.DefaultQuery("label_ids", ""), ",")

		// Array of ledger ids
		whereIn := []int{}

		// Run query
		l := []models.LabelsToLedger{}
		t.db.New().Where("LabelsToLedgerLabelId IN (?) AND LabelsToLedgerAccountId = ?", ids, accountId).Find(&l)

		// Build id array.
		for _, row := range l {
			whereIn = append(whereIn, int(row.LabelsToLedgerLedgerId))
		}

		// Update query.
		params.Wheres = append(params.Wheres, models.KeyValue{
			Key:          "LedgerId",
			Compare:      "IN",
			ValueIntList: whereIn,
		})
	}

	// Run the query
	meta, err := t.db.QueryMeta(&results, params)

	// Loop through and add signed urls to files TODO(spicer): clean this up. Maybe move all this into the model.
	for key, row := range results {
		for key2, row2 := range row.Files {
			results[key].Files[key2].Url = t.db.GetSignedFileUrl(row2.Path)
			results[key].Files[key2].Thumb600By600Url = t.db.GetSignedFileUrl(row2.ThumbPath)
		}
	}

	// Return json based on if this was a good result or not.
	response.ResultsMeta(c, results, err, meta)
}

//
// GetLedger by id
//
func (t *Controller) GetLedger(c *gin.Context) {
	// Set id
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	// Get ledger and make sure we have perms to it
	l, err := t.db.GetLedgerByAccountAndId(uint(c.MustGet("accountId").(int)), uint(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ledger entry not found."})
		return
	}

	// Return happy.
	response.Results(c, l, nil)
}

//
// CreateLedger - Create a ledger within the account.
//
func (t *Controller) CreateLedger(c *gin.Context) {
	// Setup Contact obj
	o := models.Ledger{}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &o, "create") != nil {
		return
	}

	// Make sure the AccountId is correct.
	o.AccountId = uint(c.MustGet("accountId").(int))

	// Add in auto fields
	o.AddedById = uint(c.MustGet("userId").(int))

	// Create category
	t.db.LedgerCreate(&o)

	// Return happy.
	response.RespondCreated(c, o, nil)
}

//
// UpdateLedger - Update a ledger within the account.
//
func (t *Controller) UpdateLedger(c *gin.Context) {
	// Get ID from URL
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	// Get ledger and make sure we have perms to it
	org, err := t.db.GetLedgerByAccountAndId(uint(c.MustGet("accountId").(int)), uint(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ledger entry not found."})
		return
	}

	// Setup Ledger obj
	o := models.Ledger{}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &o, "update") != nil {
		return
	}

	// Make sure the AccountId is correct.
	o.AccountId = uint(c.MustGet("accountId").(int))

	// Change some CreatedAt dates. This is hacky. Maybe find a better way some day.
	o.CreatedAt = org.CreatedAt
	o.Contact.CreatedAt = org.Contact.CreatedAt
	o.Category.CreatedAt = org.Category.CreatedAt

	// Create category
	t.db.LedgerUpdate(&o)

	// Return happy.
	response.RespondUpdated(c, o, nil)
}

//
// DeleteLedger a ledger within the account.
//
func (t *Controller) DeleteLedger(c *gin.Context) {
	// Get Id
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	// First we make sure this is an entry we have access to.
	_, err2 := t.db.GetLedgerByAccountAndId(uint(c.MustGet("accountId").(int)), uint(id))

	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ledger entry not found."})
		return
	}

	// Delete ledger
	err = t.db.DeleteLedgerByAccountAndId(uint(c.MustGet("accountId").(int)), uint(id))

	if err != nil {
		response.RespondError(c, err)
		return
	}

	// Return happy.
	response.RespondDeleted(c, nil)
}

/* End File */
