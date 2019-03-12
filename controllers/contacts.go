//
// Date: 2019-01-13
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2019-01-13
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"strings"

	"github.com/cloudmanic/skyclerk.com/library/response"
	"github.com/cloudmanic/skyclerk.com/models"
	"github.com/gin-gonic/gin"
)

//
// Create a contact within the account.
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

/* End File */
