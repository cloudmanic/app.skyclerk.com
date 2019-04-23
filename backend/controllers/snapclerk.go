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
		Amount:    amount,
		AccountId: accountId,
		AddedById: userId,
		Contact:   c.PostForm("contact"),
		Category:  c.PostForm("category"),
		Labels:    c.PostForm("labels"),
		Note:      c.PostForm("note"),
		Lat:       c.PostForm("lat"),
		Lon:       c.PostForm("lon"),
		Status:    "Pending",
		FileId:    o.Id,
		File:      o,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	}

	// Store in DB
	t.db.SnapClerkCreate(&sc)

	// Return happy.
	response.RespondCreated(c, sc, nil)
}

/* End File */
