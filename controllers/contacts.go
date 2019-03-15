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

	"github.com/cloudmanic/skyclerk.com/library/request"
	"github.com/cloudmanic/skyclerk.com/library/response"
	"github.com/cloudmanic/skyclerk.com/models"
)

//
// GetContacts - Return a list of contacts. We limit to 500 mainly so we do not overload the
// system, but enough so the front-end does not have to page
//
func (t *Controller) GetContacts(c *gin.Context) {

	// Place to store the results.
	var results = []models.Contact{}

	// Get limits and pages
	page, _, _ := request.GetSetPagingParms(c)

	// Set the query parms
	params := models.QueryParam{
		Order:            c.DefaultQuery("order", "ContactsName"),
		Sort:             c.DefaultQuery("sort", "ASC"),
		Limit:            500,
		Page:             page,
		AllowedOrderCols: []string{"ContactsId", "ContactsName"},
		Wheres: []models.KeyValue{
			{Key: "ContactsAccountId", ValueInt: c.MustGet("accountId").(int)},
		},
	}

	// Run the query
	meta, err := t.db.QueryMeta(&results, params)

	// Return json based on if this was a good result or not.
	response.ResultsMeta(c, results, err, meta)
}

//
// GetContact - Get Contact by id
//
func (t *Controller) GetContact(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	// Get contact and make sure we have perms to it
	orgCon, err := t.db.GetContactByAccountAndId(uint(c.MustGet("accountId").(int)), uint(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Contact not found."})
		return
	}

	// Return happy.
	response.Results(c, orgCon, nil)
}

//
// CreateContact - Create a contact within the account.
//
func (t *Controller) CreateContact(c *gin.Context) {

	// Setup Contact obj
	o := models.Contact{}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &o, "create") != nil {
		return
	}

	// Make sure the AccountId is correct.
	o.AccountId = uint(c.MustGet("accountId").(int))

	// Clean up some vars
	o.Name = strings.Trim(o.Name, " ")
	o.FirstName = strings.Trim(o.FirstName, " ")
	o.LastName = strings.Trim(o.LastName, " ")

	// Create category
	t.db.New().Create(&o)

	// Return happy.
	response.RespondCreated(c, o, nil)
}

//
// UpdateContact - Update a contact within the account.
//
func (t *Controller) UpdateContact(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	// First we make sure this is an entry we have access to.
	orgCon, err := t.db.GetContactByAccountAndId(uint(c.MustGet("accountId").(int)), uint(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Contact not found."})
		return
	}

	// Setup Contact obj
	o := models.Contact{}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &o, "update") != nil {
		return
	}

	// We just allow updating of a few fields
	orgCon.Name = strings.Trim(o.Name, " ")
	orgCon.FirstName = strings.Trim(o.FirstName, " ")
	orgCon.LastName = strings.Trim(o.LastName, " ")

	// Update category
	t.db.New().Save(&orgCon)

	// Return happy.
	response.RespondUpdated(c, orgCon, nil)
}

//
// DeleteContact a contact within the account.
//
func (t *Controller) DeleteContact(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	// First we make sure this is an entry we have access to.
	_, err2 := t.db.GetContactByAccountAndId(uint(c.MustGet("accountId").(int)), uint(id))

	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Contact not found."})
		return
	}

	// Delete category
	err = t.db.DeleteContactByAccountAndId(uint(c.MustGet("accountId").(int)), uint(id))

	if err != nil {
		response.RespondError(c, err)
		return
	}

	// Return happy.
	response.RespondDeleted(c, nil)
}

/* End File */
