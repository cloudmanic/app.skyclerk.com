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
// GetContacts - Return a list of contacts. We limit to 500 mainly so we do not overload the
// system, but enough so the front-end does not have to page
//
func (t *Controller) GetContacts(c *gin.Context) {
	// Place to store the results.
	var results = []models.Contact{}

	// Get the account id
	accountId := c.MustGet("accountId").(int)

	// Get limits and pages
	page, _, _ := request.GetSetPagingParms(c)

	// Get limit
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "500"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	// Set the query parms
	params := models.QueryParam{
		Order:            c.DefaultQuery("order", "ContactsName"),
		Sort:             c.DefaultQuery("sort", "ASC"),
		Limit:            limit,
		Page:             page,
		AllowedOrderCols: []string{"ContactsId", "ContactsName"},
		Wheres: []models.KeyValue{
			{Key: "ContactsAccountId", Compare: "=", ValueInt: accountId},
		},
	}

	// Manage a search query - TODO(spicer): this is hacky but it will do for now.
	if len(c.DefaultQuery("search", "")) > 0 {
		// Query term
		search := c.DefaultQuery("search", "")

		// Array of contact ids
		contactWhereIn := []int{}

		// Get contacts
		con := []models.Contact{}
		t.db.New().Where("(ContactsName LIKE ? OR ContactsFirstName LIKE ? OR ContactsLastName LIKE ?) AND ContactsAccountId = ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", accountId).Find(&con)

		// Build id array.
		for _, row := range con {
			contactWhereIn = append(contactWhereIn, int(row.Id))
		}

		// Update query.
		params.Wheres = append(params.Wheres, models.KeyValue{
			Key:          "ContactsId",
			Compare:      "IN",
			ValueIntList: contactWhereIn,
		})

		// If we have no ids we have no results
		if len(contactWhereIn) == 0 {
			// Run the query
			meta, err := t.db.QueryMeta(&results, params)

			// Return json based on if this was a good result or not.
			response.ResultsMeta(c, []models.Contact{}, err, meta)
			return
		}
	}

	// Run the query
	meta, err := t.db.QueryMeta(&results, params)

	// TODO(spicer): Move this into the model maybe.
	for key, row := range results {
		// Double check the contact has an avatar. This is just to double check.
		t.db.ConfirmContactAvatar(&row)

		// Add a signed avatar path
		results[key].AvatarUrl = t.db.GetSignedFileUrl(row.Avatar)
	}

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
	o.Email = strings.Trim(o.Email, " ")
	o.Phone = strings.Trim(o.Phone, " ")
	o.Address = strings.Trim(o.Address, " ")
	o.City = strings.Trim(o.City, " ")
	o.State = strings.Trim(o.State, " ")
	o.Zip = strings.Trim(o.Zip, " ")
	o.Country = strings.Trim(o.Country, " ")
	o.Twitter = strings.Trim(o.Twitter, " ")
	o.Facebook = strings.Trim(o.Facebook, " ")
	o.Linkedin = strings.Trim(o.Linkedin, " ")
	o.Website = strings.Trim(o.Website, " ")

	// Create category
	err := t.db.CreateContact(&o)

	if err != nil {
		services.Critical(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong creating your contact. Please contact help@skyclerk.com"})
		return
	}

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
	orgCon.Email = strings.Trim(o.Email, " ")
	orgCon.Phone = strings.Trim(o.Phone, " ")
	orgCon.Address = strings.Trim(o.Address, " ")
	orgCon.City = strings.Trim(o.City, " ")
	orgCon.State = strings.Trim(o.State, " ")
	orgCon.Zip = strings.Trim(o.Zip, " ")
	orgCon.Country = strings.Trim(o.Country, " ")
	orgCon.Twitter = strings.Trim(o.Twitter, " ")
	orgCon.Facebook = strings.Trim(o.Facebook, " ")
	orgCon.Linkedin = strings.Trim(o.Linkedin, " ")
	orgCon.Website = strings.Trim(o.Website, " ")

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
