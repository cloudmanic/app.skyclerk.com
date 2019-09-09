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

	"app.skyclerk.com/backend/library/helpers"
	"app.skyclerk.com/backend/library/request"
	"app.skyclerk.com/backend/library/response"
	"app.skyclerk.com/backend/models"
	"app.skyclerk.com/backend/services"
)

// Summary results
type LedgerSummary struct {
	Years      []LedgerYearSummaryResult `json:"years"`
	Labels     []LedgerSummaryResult     `json:"labels"`
	Categories []LedgerSummaryResult     `json:"categories"`
}

// LedgerSummaryResult
type LedgerSummaryResult struct {
	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// LedgerYearSummaryResult
type LedgerYearSummaryResult struct {
	Year  int `json:"year"`
	Count int `json:"count"`
}

//
// GetLedgers - Return a list of ledgers. We limit to 25 mainly so we do not overload the
// system, but enough so the front-end does not have to page
//
func (t *Controller) GetLedgers(c *gin.Context) {
	// Query database based on url parms.
	results, meta, err := t.QueryLedgers(c, 25, []string{"Category", "Contact", "Labels", "Files"})

	// Error responses were already set in QueryLedgers
	if err != nil {
		return
	}

	// Loop through and add signed urls to files TODO(spicer): clean this up. Maybe move all this into the model.
	for key, row := range results {
		for key2, row2 := range row.Files {
			results[key].Files[key2].Url = t.db.GetSignedFileUrl(row2.Path)
			results[key].Files[key2].Thumb600By600Url = t.db.GetSignedFileUrl(row2.ThumbPath)
		}
	}

	// TODO(spicer): Move this into the model maybe.
	for key, row := range results {
		// Double check the contact has an avatar. This is just to double check.
		t.db.ConfirmContactAvatar(&row.Contact)

		// Add a signed avatar path
		results[key].Contact.AvatarUrl = t.db.GetSignedFileUrl(row.Contact.Avatar)
	}

	// Return json based on if this was a good result or not.
	response.ResultsMeta(c, results, err, meta)
}

//
// GetLedger by id
//
func (t *Controller) GetLedger(c *gin.Context) {
	// Set id
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	// Get ledger and make sure we have perms to it
	l, err := t.db.GetLedgerByAccountAndId(uint(c.MustGet("accountId").(int)), uint(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ledger entry not found."})
		return
	}

	// Double check the contact has an avatar. This is just to double check.
	t.db.ConfirmContactAvatar(&l.Contact)

	// Add a signed avatar path
	l.Contact.AvatarUrl = t.db.GetSignedFileUrl(l.Contact.Avatar)

	// Return happy.
	response.Results(c, l, nil)
}

//
// CreateLedger - Create a ledger within the account.
//
func (t *Controller) CreateLedger(c *gin.Context) {
	// Setup Contact obj
	o := models.Ledger{}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &o, "create") != nil {
		return
	}

	// Make sure the AccountId is correct.
	o.AccountId = uint(c.MustGet("accountId").(int))

	// Add in auto fields
	o.AddedById = uint(c.MustGet("userId").(int))

	// Create ledger
	t.db.LedgerCreate(&o)

	// Fresh pull
	j, err := t.db.GetLedgerByAccountAndId(o.AccountId, o.Id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "System error. Please contact help@skyclerk.com."})
		return
	}

	// Set the ledger type
	ledgerType := "expense"

	if o.Amount > 0 {
		ledgerType = "income"
	}

	// Get the contact name.
	contactName := o.Contact.Name

	if len(contactName) == 0 {
		contactName = o.Contact.FirstName + " " + o.Contact.LastName
	}

	// Add to the activity log
	t.db.New().Create(&models.Activity{
		AccountId: o.AccountId,
		UserId:    o.AddedById,
		Action:    ledgerType,
		SubAction: "create",
		Name:      contactName,
		Amount:    o.Amount,
		LedgerId:  o.Id,
	})

	// Return happy.
	response.RespondCreated(c, j, nil)
}

//
// UpdateLedger - Update a ledger within the account.
//
func (t *Controller) UpdateLedger(c *gin.Context) {
	// Get ID from URL
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	// Get ledger and make sure we have perms to it
	org, err := t.db.GetLedgerByAccountAndId(uint(c.MustGet("accountId").(int)), uint(id))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ledger entry not found."})
		return
	}

	// Setup Ledger obj
	o := models.Ledger{}

	// Here we parse the JSON sent in, assign it to a struct, set validation errors if any.
	if t.ValidateRequest(c, &o, "update") != nil {
		return
	}

	// Make sure the AccountId is correct.
	o.AccountId = uint(c.MustGet("accountId").(int))

	// Change some CreatedAt dates. This is hacky. Maybe find a better way some day.
	o.CreatedAt = org.CreatedAt
	o.Contact.CreatedAt = org.Contact.CreatedAt
	o.Category.CreatedAt = org.Category.CreatedAt

	// Create category
	t.db.LedgerUpdate(&o)

	// Set the ledger type
	ledgerType := "expense"

	if o.Amount > 0 {
		ledgerType = "income"
	}

	// Get the contact name.
	contactName := o.Contact.Name

	if len(contactName) == 0 {
		contactName = o.Contact.FirstName + " " + o.Contact.LastName
	}

	// Add to the activity log
	t.db.New().Create(&models.Activity{
		AccountId: o.AccountId,
		UserId:    uint(c.MustGet("userId").(int)),
		Action:    ledgerType,
		SubAction: "update",
		Name:      contactName,
		Amount:    o.Amount,
		LedgerId:  o.Id,
	})

	// Return happy.
	response.RespondUpdated(c, o, nil)
}

//
// DeleteLedger a ledger within the account.
//
func (t *Controller) DeleteLedger(c *gin.Context) {
	// Get Id
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": err})
		return
	}

	// First we make sure this is an entry we have access to.
	entry, err2 := t.db.GetLedgerByAccountAndId(uint(c.MustGet("accountId").(int)), uint(id))

	if err2 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ledger entry not found."})
		return
	}

	// Delete ledger
	err = t.db.DeleteLedgerByAccountAndId(uint(c.MustGet("accountId").(int)), uint(id))

	if err != nil {
		response.RespondError(c, err)
		return
	}

	// Set the ledger type
	ledgerType := "expense"

	if entry.Amount > 0 {
		ledgerType = "income"
	}

	// Get the contact name.
	contactName := entry.Contact.Name

	if len(contactName) == 0 {
		contactName = entry.Contact.FirstName + " " + entry.Contact.LastName
	}

	// Add to the activity log
	t.db.New().Create(&models.Activity{
		AccountId: entry.AccountId,
		UserId:    uint(c.MustGet("userId").(int)),
		Action:    ledgerType,
		SubAction: "delete",
		Name:      contactName,
		Amount:    entry.Amount,
		LedgerId:  entry.Id,
	})

	// Return happy.
	response.RespondDeleted(c, nil)
}

//
// GetLedgerPlSummary return the P&L for the ledger via parms passed in.
//
func (t *Controller) GetLedgerPlSummary(c *gin.Context) {
	// Values we track.
	income := 0.00
	expense := 0.00

	// Query database based on url parms.
	results, _, err := t.QueryLedgers(c, 0, []string{})

	// Error responses were already set in QueryLedgers
	if err != nil {
		return
	}

	// Add up the values we track
	for _, row := range results {
		if row.Amount > 0 {
			income = income + row.Amount
		} else {
			expense = expense + (row.Amount * -1)
		}
	}

	// set profit
	profit := (income - expense)

	//Return success json.
	c.JSON(200, gin.H{"income": helpers.Round(income, 2), "expense": helpers.Round((expense * -1), 2), "profit": helpers.Round(profit, 2)})
}

//
// GetLedgerSummary get a smmary of the ledger
//
func (t *Controller) GetLedgerSummary(c *gin.Context) {
	// Get account
	accountId := c.MustGet("accountId").(int)

	// Setup return object
	ls := LedgerSummary{}

	// Build SQL for category (Notice: lower case column names)
	catSql := "SELECT CategoriesId AS id, CategoriesName AS name, COUNT(CategoriesId) AS count FROM Ledger INNER JOIN Categories ON Categories.CategoriesId = Ledger.LedgerCategoryId WHERE (LedgerAccountId = ?)"

	// Build SQL for labels (Notice: lower case column names)
	lbsSql := "SELECT LabelsToLedgerLabelId AS id, LabelsName AS name, COUNT(LabelsToLedgerLabelId) AS count FROM `LabelsToLedger` INNER JOIN `Ledger` ON `LabelsToLedger`.`LabelsToLedgerLedgerId` = `Ledger`.`LedgerId` INNER JOIN `Labels` ON `LabelsToLedger`.`LabelsToLedgerLabelId` = `Labels`.`LabelsId` WHERE (LedgerAccountId = ?)"

	// Build SQL for years (Notice: lower case column names)
	yrsSql := "SELECT YEAR(LedgerDate) as year, COUNT(LedgerId) as count FROM Ledger WHERE (LedgerAccountId = ?)"

	// Add type filter - income
	if c.DefaultQuery("type", "") == "income" {
		catSql = catSql + " AND (LedgerAmount >= 0.00)"
		lbsSql = lbsSql + " AND (LedgerAmount >= 0.00)"
		yrsSql = yrsSql + " AND (LedgerAmount >= 0.00)"
	}

	// Add type filter - expense
	if c.DefaultQuery("type", "") == "expense" {
		catSql = catSql + " AND (LedgerAmount < 0.00)"
		lbsSql = lbsSql + " AND (LedgerAmount < 0.00)"
		yrsSql = yrsSql + " AND (LedgerAmount < 0.00)"
	}

	// Add Group
	catSql = catSql + " GROUP BY CategoriesName ORDER BY CategoriesName ASC"
	lbsSql = lbsSql + " GROUP BY LabelsToLedgerLabelId ORDER BY LabelsName ASC"
	yrsSql = yrsSql + " GROUP BY YEAR(LedgerDate) ORDER BY LedgerDate DESC"

	// Run query.
	t.db.New().Raw(yrsSql, accountId).Scan(&ls.Years)
	t.db.New().Raw(lbsSql, accountId).Scan(&ls.Labels)
	t.db.New().Raw(catSql, accountId).Scan(&ls.Categories)

	// Return happy.
	response.Results(c, ls, nil)
}

// -------------- Private Helper Functions ------------------ //

//
// QueryLedgers is shared between a few functions so we break it out.
// This function is not a route function. We pass in preloads to make for
// less queries if we don't need them.
//
func (t *Controller) QueryLedgers(c *gin.Context, limit int, preloads []string) ([]models.Ledger, models.QueryMetaData, error) {
	// Place to store the results.
	var results = []models.Ledger{}

	// Get account
	accountId := c.MustGet("accountId").(int)

	// Get limits and pages
	page, _, _ := request.GetSetPagingParms(c)

	// Set the query parms
	params := models.QueryParam{
		Order:            c.DefaultQuery("order", "LedgerDate"),
		Sort:             c.DefaultQuery("sort", "DESC"),
		Page:             page,
		Debug:            false,
		PreLoads:         preloads,
		AllowedOrderCols: []string{"LedgerId", "LedgerDate"},
		Wheres: []models.KeyValue{
			{Key: "LedgerAccountId", Compare: "=", ValueInt: accountId},
		},
	}

	// Set limit
	if limit > 0 {
		params.Limit = limit
	}

	// Add type filter - income
	if c.DefaultQuery("type", "") == "income" {
		params.Wheres = append(params.Wheres, models.KeyValue{
			Key:        "LedgerAmount",
			Compare:    ">=",
			ValueFloat: 0.01, // because of the lib we can't use 0
		})
	}

	// Add type filter - expense
	if c.DefaultQuery("type", "") == "expense" {
		params.Wheres = append(params.Wheres, models.KeyValue{
			Key:        "LedgerAmount",
			Compare:    "<=",
			ValueFloat: -0.01, // because of the lib we can't use 0
		})
	}

	// Add type filter - category_id
	if len(c.DefaultQuery("category_id", "")) > 0 {
		// Convert cat id
		cat_id, err := strconv.Atoi(c.DefaultQuery("category_id", ""))

		if err != nil {
			services.Info(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error with category_id"})
			return results, models.QueryMetaData{}, err
		}

		// Update query.
		params.Wheres = append(params.Wheres, models.KeyValue{
			Key:      "LedgerCategoryId",
			Compare:  "=",
			ValueInt: cat_id,
		})
	}

	// Add type filter - year
	if len(c.DefaultQuery("year", "")) > 0 {
		params.Wheres = append(params.Wheres, models.KeyValue{
			Key:     "LedgerDate",
			Compare: ">=",
			Value:   c.DefaultQuery("year", "") + "-01-01",
		})

		params.Wheres = append(params.Wheres, models.KeyValue{
			Key:     "LedgerDate",
			Compare: "<=",
			Value:   c.DefaultQuery("year", "") + "-12-31",
		})
	}

	// Add type filter - start date
	if len(c.DefaultQuery("start_date", "")) > 0 {
		params.Wheres = append(params.Wheres, models.KeyValue{
			Key:     "LedgerDate",
			Compare: ">=",
			Value:   c.DefaultQuery("start_date", ""),
		})
	}

	// Add type filter - end date
	if len(c.DefaultQuery("end_date", "")) > 0 {
		params.Wheres = append(params.Wheres, models.KeyValue{
			Key:     "LedgerDate",
			Compare: "<=",
			Value:   c.DefaultQuery("end_date", ""),
		})
	}

	// Get ledger ids from lables we want to filter from.
	if len(c.DefaultQuery("label_ids", "")) > 0 {
		// Get Ids from url
		ids := strings.Split(c.DefaultQuery("label_ids", ""), ",")

		// Array of ledger ids
		whereIn := []int{}

		// Run query
		l := []models.LabelsToLedger{}
		t.db.New().Where("LabelsToLedgerLabelId IN (?)", ids).Find(&l)

		// Build id array.
		for _, row := range l {
			whereIn = append(whereIn, int(row.LabelsToLedgerLedgerId))
		}

		// Update query.
		params.Wheres = append(params.Wheres, models.KeyValue{
			Key:          "LedgerId",
			Compare:      "IN",
			ValueIntList: whereIn,
		})
	}

	// Manage a search query
	if len(c.DefaultQuery("search", "")) > 0 {
		// Query term
		search := c.DefaultQuery("search", "")

		// Array of contact ids
		contactWhereIn := []int{}

		// Get contacts
		c := []models.Contact{}
		t.db.New().Where("(ContactsName LIKE ? OR ContactsFirstName LIKE ? OR ContactsLastName LIKE ?) AND ContactsAccountId = ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", accountId).Find(&c)

		// Build id array.
		for _, row := range c {
			contactWhereIn = append(contactWhereIn, int(row.Id))
		}

		// Update query.
		params.Wheres = append(params.Wheres, models.KeyValue{
			Key:          "LedgerContactId",
			Compare:      "IN",
			ValueIntList: contactWhereIn,
		})
	}

	// Run the query
	meta, err := t.db.QueryMeta(&results, params)

	if err != nil {
		services.Info(err)
		return results, meta, err
	}

	// Return happy
	return results, meta, nil
}

/* End File */
