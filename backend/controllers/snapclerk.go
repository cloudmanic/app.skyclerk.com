//
// Date: 2019-04-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package controllers

import (
	"errors"
	"flag"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"app.skyclerk.com/backend/emails"
	"app.skyclerk.com/backend/library/email"
	"app.skyclerk.com/backend/library/request"
	"app.skyclerk.com/backend/library/response"
	"app.skyclerk.com/backend/library/store/object"
	"app.skyclerk.com/backend/models"
	"app.skyclerk.com/backend/services"
)

//
// GetSnapClerkUsage - returns the monthly usage.
//
func (t *Controller) GetSnapClerkUsage(c *gin.Context) {
	// AccountId.
	accountId := uint(c.MustGet("accountId").(int))

	// Get the total number of SnapClerks used.
	count := t.db.SnapClerkMonthlyUsage(accountId)

	type r struct {
		Count int `json:"count"`
	}

	rt := r{Count: count}

	// Return happy.
	response.Results(c, rt, nil)
}

//
// CreateSnapClerkByFileId - Create a new SnapClerk by file_id
//
func (t *Controller) CreateSnapClerkByFileId(c *gin.Context) {
	// UserId.
	userID := uint(c.MustGet("userId").(int))

	// AccountId.
	accountID := uint(c.MustGet("accountId").(int))

	// Setup Contact obj
	o := models.SnapClerk{}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &o, "create") != nil {
		return
	}

	// Make sure the AccountId is correct.
	o.AddedById = userID
	o.AccountId = accountID

	// Store in DB
	t.db.SnapClerkCreate(&o)

	// Refresh object
	sc, err := t.db.GetSnapClerkByAccountAndId(accountID, o.Id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Snapclerk not found."})
		return
	}

	// Add to AppLog
	t.CreateActivityLogEntry(sc)

	// Notify users we received this.
	t.NoifyReceiptWasReceived(sc)

	// Return happy.
	response.RespondCreated(c, sc, nil)
}

//
// CreateSnapClerk - Upload a file to store in snapclerk
//
func (t *Controller) CreateSnapClerk(c *gin.Context) {
	// UserId.
	userID := uint(c.MustGet("userId").(int))

	// AccountId.
	accountID := uint(c.MustGet("accountId").(int))

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
		AccountId:    accountID,
		AddedById:    userID,
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

	// Add to AppLog
	t.CreateActivityLogEntry(sc)

	// Notify users we received this.
	t.NoifyReceiptWasReceived(sc)

	// Return happy.
	response.RespondCreated(c, sc, nil)
}

//
// GetSnapClerk - Return a list of snapclerk. We limit to 100 mainly so we do not overload the
// system, but enough so the front-end does not have to page
//
func (t *Controller) GetSnapClerk(c *gin.Context) {
	// Place to store the results.
	var results = []models.SnapClerk{}

	// Get limits and pages
	page, limit, _ := request.GetSetPagingParms(c)

	// Set the query parms
	params := models.QueryParam{
		Order:            c.DefaultQuery("order", "SnapClerkId"),
		Sort:             c.DefaultQuery("sort", "ASC"),
		Limit:            limit,
		Page:             page,
		Debug:            false,
		PreLoads:         []string{"File"},
		AllowedOrderCols: []string{"SnapClerkId"},
		Wheres: []models.KeyValue{
			{Key: "SnapClerkAccountId", Compare: "=", ValueInt: c.MustGet("accountId").(int)},
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

//
// CreateActivityLogEntry records a entry in the app long table for the snapclerk that was uploaded.
//
func (t *Controller) CreateActivityLogEntry(snapClerk models.SnapClerk) {
	// Add to the activity log
	t.db.New().Create(&models.Activity{
		AccountId:   snapClerk.AccountId,
		UserId:      snapClerk.AddedById,
		Action:      "snapclerk",
		SubAction:   "create",
		SnapClerkId: snapClerk.Id,
	})
}

//
// NoifyReceiptWasReceived send notices that we received.
//
// We send emails in the controller instead of the model as it seems
// there might be times we want to do work within the models and not
// send email. We assume the only time we need to send email is when the call
// came via an API call.
//
// Also we do it in the controller because of circular import issues.
//
func (t *Controller) NoifyReceiptWasReceived(snapClerk models.SnapClerk) {
	// If we do not have a file send an error.
	if len(snapClerk.File.Path) == 0 {
		services.Info(errors.New("no file was passed into NoifyReceiptWasReceived()"))
	}

	// Make sure we have an account id.
	if snapClerk.AccountId <= 0 {
		services.Info(errors.New("no AccountId was passed into NoifyReceiptWasReceived()"))
		return
	}

	// Download file so we can cache it locally.
	filePath, err := object.DownloadObject(snapClerk.File.Path)

	if err != nil {
		services.Info(err)
	}

	// Get all users on this account.
	users := t.db.GetUsersByAccount(snapClerk.AccountId)

	// Send notices to users that that receipt was received.
	for _, row := range users {
		sendSnapClerkReceipt(row, filePath)
	}
}

//
// sendSnapClerkReceipt - send an email with receipt of snapclerk
//
func sendSnapClerkReceipt(user models.User, filePath string) {
	// Set subject.
	subject := "We've Received Your Snap!Clerk Receipt"

	// Build attachments.
	attachments := []string{filePath}

	// Send welcome email to user already in the system.
	if flag.Lookup("test.v") != nil {
		email.Send(user.Email, "", subject, emails.GetSnapClerkReceiptHTML(user), attachments)
	} else {
		go email.Send(user.Email, "", subject, emails.GetSnapClerkReceiptHTML(user), attachments)
	}
}

/* End File */
