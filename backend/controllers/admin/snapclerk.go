//
// Date: 11/3/2018
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package admin

import (
	"errors"
	"flag"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"optionsnews.com/services"

	"app.skyclerk.com/backend/emails"
	"app.skyclerk.com/backend/library/email"
	"app.skyclerk.com/backend/library/helpers"
	"app.skyclerk.com/backend/library/store/object"
	"app.skyclerk.com/backend/models"
)

//
// GetSnapClerks - Get a list of snap clerks.
//
func (t *Controller) GetSnapClerks(c *gin.Context) {
	// SnapClerk to return
	results := []models.SnapClerk{}

	// Make query
	t.db.New().Preload("File").Where("SnapClerkStatus = ?", "Pending").Find(&results)

	// Loop through and add signed urls to files TODO(spicer): clean this up. Maybe move all this into the model.
	for key, row := range results {
		results[key].File.Url = t.db.GetSignedFileUrl(row.File.Path)
		results[key].File.Thumb600By600Url = t.db.GetSignedFileUrl(row.File.ThumbPath)
	}

	// Return happy JSON
	c.JSON(200, results)
}

//
// ConvertSnapClerk - convert a snapclerk to a ledger entry
//
func (t *Controller) ConvertSnapClerk(c *gin.Context) {
	// Get SnapClerk Id
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// Read the JSON POSTed in.
	body, _ := ioutil.ReadAll(c.Request.Body)
	amount := gjson.Get(string(body), "amount").Float()
	accountID := gjson.Get(string(body), "account_id").Int()
	contact := gjson.Get(string(body), "contact").String()
	category := gjson.Get(string(body), "category").String()
	createdAt := gjson.Get(string(body), "created_at").String()
	note := gjson.Get(string(body), "note").String()

	// Get OG snapclerk
	og := models.SnapClerk{}
	t.db.New().Find(&og, id)

	// Get snapclerk by ID
	sc, err := t.db.GetSnapClerkByAccountAndId(og.AccountId, uint(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// Store created at date.
	uploadDate := sc.CreatedAt

	// Update Snap!Clerk with the new values
	sc.Amount = (amount * -1)
	sc.Note = strings.Trim(note, " ")
	sc.Contact = strings.Trim(contact, " ")
	sc.Category = strings.Trim(category, " ")
	sc.AccountId = uint(accountID)

	// Add in date from CreatedAt
	sc.CreatedAt = helpers.ParseDateNoError(createdAt)

	// Convert a snapclerk to a ledger entry.
	ledger, err := t.db.ConvertSnapclerkToLedger(sc)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// Add in the ledger id, mark as done and save
	sc.CreatedAt = uploadDate
	sc.Status = "Processed"
	sc.LedgerId = ledger.Id
	sc.ProcessedAt = time.Now()
	sc.ReviewedById = uint(c.MustGet("userId").(int))
	t.db.New().Save(&sc)

	// Fresh lookup
	l, err := t.db.GetLedgerByAccountAndId(uint(accountID), ledger.Id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// Notify users the receipt has been Processed
	t.NoifyReceiptWasProcessed(sc, l)

	// Return happy JSON
	c.JSON(200, l)
}

//
// RejectSnapClerk - reject a snapclerk
//
func (t *Controller) RejectSnapClerk(c *gin.Context) {
	// Get SnapClerk Id
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// Get OG snapclerk
	og := models.SnapClerk{}
	t.db.New().Find(&og, id)

	// Get snapclerk by ID
	sc, err := t.db.GetSnapClerkByAccountAndId(og.AccountId, uint(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// Mark rejected
	sc.Status = "Rejected"
	sc.ProcessedAt = time.Now()
	sc.ReviewedById = uint(c.MustGet("userId").(int))
	t.db.New().Save(&sc)

	// Notify users the receipt has been Rejected
	t.NoifyReceiptWasRejected(sc)

	// Return happy JSON
	c.JSON(204, "")
}

//
// NoifyReceiptWasProcessed send notices that we Processed.
//
//
func (t *Controller) NoifyReceiptWasProcessed(snapClerk models.SnapClerk, ledger models.Ledger) {
	// Make sure we have an account id.
	if snapClerk.AccountId <= 0 {
		services.Info(errors.New("no AccountId was passed into NoifyReceiptWasProcessed()"))
		return
	}
	// Get all users on this account.
	users := t.db.GetUsersByAccount(snapClerk.AccountId)

	// Send notices to users that that receipt was received.
	for _, row := range users {
		sendSnapClerkProcessedEmail(row, ledger)
	}
}

//
// sendSnapClerkProcessedEmail - send an email with receipt of snapclerk
//
func sendSnapClerkProcessedEmail(user models.User, ledger models.Ledger) {
	// Set subject.
	subject := "Your Snap!Clerk Receipt Has Been Processed"

	// Build attachments.
	attachments := []string{}

	// Send email that we Processed the receipt.
	if flag.Lookup("test.v") != nil {
		email.Send(user.Email, subject, emails.GetSnapClerkProcessedHTML(user, ledger), attachments)
	} else {
		go email.Send(user.Email, subject, emails.GetSnapClerkProcessedHTML(user, ledger), attachments)
	}
}

//
// NoifyReceiptWasRejected send notices that we Rejected.
//
//
func (t *Controller) NoifyReceiptWasRejected(snapClerk models.SnapClerk) {
	// If we do not have a file send an error.
	if len(snapClerk.File.Path) == 0 {
		services.Info(errors.New("no file was passed into NoifyReceiptWasRejected()"))
	}

	// Make sure we have an account id.
	if snapClerk.AccountId <= 0 {
		services.Info(errors.New("no AccountId was passed into NoifyReceiptWasRejected()"))
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
		sendSnapClerkRejectedEmail(row, filePath)
	}
}

//
// sendSnapClerkRejectedEmail - send an email with Rejected of snapclerk
//
func sendSnapClerkRejectedEmail(user models.User, filePath string) {
	// Set subject.
	subject := "Your Snap!Clerk Receipt Was Rejected"

	// Build attachments.
	attachments := []string{filePath}

	// Send email that we Processed the receipt.
	if flag.Lookup("test.v") != nil {
		email.Send(user.Email, subject, emails.GetSnapClerkRejectedHTML(user), attachments)
	} else {
		go email.Send(user.Email, subject, emails.GetSnapClerkRejectedHTML(user), attachments)
	}
}

/* End File */
