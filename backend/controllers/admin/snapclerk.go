//
// Date: 11/3/2018
// Author(s): Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package admin

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"

	"app.skyclerk.com/backend/library/helpers"
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
	sc.Amount = amount
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
	t.db.New().Save(&sc)

	// Fresh lookup
	l, err := t.db.GetLedgerByAccountAndId(uint(accountID), ledger.Id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}

	// Return happy JSON
	c.JSON(200, l)
}

/* End File */
