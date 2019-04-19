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

	"app.skyclerk.com/backend/library/request"
	"app.skyclerk.com/backend/library/response"
	"app.skyclerk.com/backend/models"
)

//
// GetLabels - Return a list of labels. We limit to 500 mainly so we do not overload the
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
			{Key: "LabelsAccountId", ValueInt: c.MustGet("accountId").(int)},
		},
	}

	// Run the query
	meta, err := t.db.QueryMeta(&results, params)

	// Return json based on if this was a good result or not.
	response.ResultsMeta(c, results, err, meta)
}

//
// GetLabel - Get Label by id
//
func (t *Controller) GetLabel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	// Get category and make sure we have perms to it
	l, err := t.db.GetLabelByAccountAndId(uint(c.MustGet("accountId").(int)), uint(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Label not found."})
		return
	}

	// Return happy.
	response.Results(c, l, nil)
}

//
// CreateLabel - Create a Label within the account.
//
func (t *Controller) CreateLabel(c *gin.Context) {

	// Setup Label obj
	o := models.Label{}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &o, "create") != nil {
		return
	}

	// Make sure the AccountId is correct.
	o.AccountId = uint(c.MustGet("accountId").(int))

	// Clean up some vars
	o.Name = strings.Trim(o.Name, " ")

	// Create label
	t.db.New().Create(&o)

	// Return happy.
	response.RespondCreated(c, o, nil)
}

//
// UpdateLabel - Pass in a label to update.
//
func (t *Controller) UpdateLabel(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	// First we make sure this is an entry we have access to.
	orgLb, err := t.db.GetLabelByAccountAndId(uint(c.MustGet("accountId").(int)), uint(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Label not found."})
		return
	}

	// Setup Label obj
	o := models.Label{}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &o, "update") != nil {
		return
	}

	// We just allow updating of a few fields
	orgLb.Name = strings.Trim(o.Name, " ")

	// Update category
	t.db.New().Save(&orgLb)

	// Return happy.
	response.RespondUpdated(c, orgLb, nil)
}

//
// DeleteLabel - Delete a label within the account.
//
func (t *Controller) DeleteLabel(c *gin.Context) {
	// Get the label id
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	// AccountId.
	accountId := uint(c.MustGet("accountId").(int))

	// First we make sure this is an entry we have access to.
	_, err = t.db.GetLabelByAccountAndId(accountId, uint(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Label not found."})
		return
	}

	// Delete label
	err = t.db.DeleteLabelByAccountAndId(accountId, uint(id))

	if err != nil {
		response.RespondError(c, err)
		return
	}

	// Return happy.
	response.RespondDeleted(c, nil)
}

/* End File */
