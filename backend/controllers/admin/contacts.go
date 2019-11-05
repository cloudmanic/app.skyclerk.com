//
// Date: 11/3/2018
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package admin

import (
	"net/http"
	"strconv"

	"app.skyclerk.com/backend/library/helpers"
	"app.skyclerk.com/backend/library/request"
	"app.skyclerk.com/backend/library/response"
	"app.skyclerk.com/backend/models"
	"github.com/gin-gonic/gin"
)

//
// GetContacts - Return a list of contacts. We limit to 500 mainly so we do not overload the
// system, but enough so the front-end does not have to page
//
func (t *Controller) GetContacts(c *gin.Context) {
	// Place to store the results.
	var results = []models.Contact{}

	// Get the account id
	accountID := helpers.StringToInt(c.DefaultQuery("account_id", "0"))

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
			{Key: "ContactsAccountId", Compare: "=", ValueInt: accountID},
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
		t.db.New().Where("(ContactsName LIKE ? OR ContactsFirstName LIKE ? OR ContactsLastName LIKE ?) AND ContactsAccountId = ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", accountID).Find(&con)

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

/* End File */
