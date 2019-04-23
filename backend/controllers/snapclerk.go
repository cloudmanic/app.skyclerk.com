//
// Date: 2019-04-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"app.skyclerk.com/backend/library/request"
	"app.skyclerk.com/backend/library/response"
	"app.skyclerk.com/backend/models"
)

//
// CreateSnapClerk - Upload a file to store in snapclerk
//
func (t *Controller) CreateSnapClerk(c *gin.Context) {
	// UserId.
	userId := uint(c.MustGet("userId").(int))

	// AccountId.
	accountId := uint(c.MustGet("accountId").(int))

	// Do a file upload and return a file model object. Errors
	// are written to the response within this function.
	// Because of this if we have errors we simply return.
	o, err := t.DoFileUpload(c)

	if err != nil {
		return
	}

	// Convert to float (defaults to zero if not included)
	amount, _ := strconv.ParseFloat(c.PostForm("amount"), 64)

	// Build skyclerk obj from optional fields.
	sc := models.SnapClerk{
		Amount:       amount,
		AccountId:    accountId,
		AddedById:    userId,
		Contact:      c.PostForm("contact"),
		Category:     c.PostForm("category"),
		Labels:       c.PostForm("labels"),
		Note:         c.PostForm("note"),
		Lat:          c.PostForm("lat"),
		Lon:          c.PostForm("lon"),
		Status:       "Pending",
		FileId:       o.Id,
		File:         o,
		LedgerId:     0,
		ReviewedById: 0,
		UpdatedAt:    time.Now(),
		CreatedAt:    time.Now(),
	}

	// Store in DB
	t.db.SnapClerkCreate(&sc)

	// Return happy.
	response.RespondCreated(c, sc, nil)
}

//
// GetSnapClerk - Return a list of snapclerk. We limit to 50 mainly so we do not overload the
// system, but enough so the front-end does not have to page
//
func (t *Controller) GetSnapClerk(c *gin.Context) {
	// Place to store the results.
	var results = []models.SnapClerk{}

	// Get limits and pages
	page, _, _ := request.GetSetPagingParms(c)

	// Set the query parms
	params := models.QueryParam{
		Order:            c.DefaultQuery("order", "SnapClerkId"),
		Sort:             c.DefaultQuery("sort", "ASC"),
		Limit:            50,
		Page:             page,
		Debug:            false,
		PreLoads:         []string{"File"},
		AllowedOrderCols: []string{"SnapClerkId"},
		Wheres: []models.KeyValue{
			{Key: "SnapClerkAccountId", ValueInt: c.MustGet("accountId").(int)},
		},
	}

	// Run the query
	meta, err := t.db.QueryMeta(&results, params)

	// Loop through and add signed urls to files TODO(spicer): clean this up. Maybe move all this into the model.
	for key, row := range results {
		results[key].File.Url = t.db.GetSignedFileUrl(row.File.Path)
		results[key].File.Thumb600By600Url = t.db.GetSignedFileUrl(row.File.ThumbPath)
	}

	// Return json based on if this was a good result or not.
	response.ResultsMeta(c, results, err, meta)
}

/* End File */
