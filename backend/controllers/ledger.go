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

	"github.com/gin-gonic/gin"

	"github.com/cloudmanic/skyclerk.com/backend/library/request"
	"github.com/cloudmanic/skyclerk.com/backend/library/response"
	"github.com/cloudmanic/skyclerk.com/backend/models"
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

//
// GetLedger by id
//
func (t *Controller) GetLedger(c *gin.Context) {

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
