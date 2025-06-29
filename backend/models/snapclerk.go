//
// Date: 2019-04-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"

	"app.skyclerk.com/backend/library/slack"
	"app.skyclerk.com/backend/services"
)

// SnapClerk struct
type SnapClerk struct {
	Id           uint      `gorm:"primary_key;column:SnapClerkId" json:"id"`
	AccountId    uint      `gorm:"column:SnapClerkAccountId" sql:"not null" json:"account_id"`
	UpdatedAt    time.Time `gorm:"column:SnapClerkUpdatedAt" sql:"not null" json:"-"`
	CreatedAt    time.Time `gorm:"column:SnapClerkCreatedAt" sql:"not null" json:"created_at"`
	AddedById    uint      `gorm:"column:SnapClerkAddedById" sql:"not null" json:"added_by_id"`
	ReviewedById uint      `gorm:"column:SnapClerkReviewedById" sql:"not null" json:"-"`
	Status       string    `gorm:"column:SnapClerkStatus" sql:"not null;default:'Pending'" json:"status"`
	FileId       uint      `gorm:"column:SnapClerkFileId" sql:"not null" json:"file_id"`
	File         File      `gorm:"foreignkey:SnapClerkFileId" json:"file"`
	LedgerId     uint      `gorm:"column:SnapClerkLedgerId;index:SnapClerkLedgerId" sql:"not null" json:"ledger_id"`
	Amount       float64   `gorm:"column:SnapClerkAmount" sql:"not null;type:DECIMAL(12,2)" json:"amount"`
	Contact      string    `gorm:"column:SnapClerkContact" sql:"not null" json:"contact"`
	Category     string    `gorm:"column:SnapClerkCategory" sql:"not null" json:"category"`
	Labels       string    `gorm:"column:SnapClerkLabels" sql:"not null" json:"labels"`
	Note         string    `gorm:"column:SnapClerkNote" sql:"not null;type:TEXT" json:"note"`
	Lat          string    `gorm:"column:SnapClerkLat" sql:"not null;type:TEXT" json:"lat"`
	Lon          string    `gorm:"column:SnapClerkLon" sql:"not null;type:TEXT" json:"lon"`
	ProcessedAt  time.Time `gorm:"column:SnapClerkProcessedAt" sql:"not null" json:"processed_at"` // This has to be at the end GORM (auto update stuff)
}

//
// Set the table name.
//
func (SnapClerk) TableName() string {
	return "SnapClerk"
}

//
// Validate for this model.
//
func (a SnapClerk) Validate(db Datastore, action string, userId uint, accountId uint, objId uint) error {

	return validation.ValidateStruct(&a,

		validation.Field(&a.FileId,
			validation.Required.Error("The file_id field is required."),
		),
	)
}

//
// SnapClerkCreate - Create a new SnapClerk entry.
//
func (db *DB) SnapClerkCreate(sc *SnapClerk) error {
	// Store this entry.
	db.Create(&sc)

	// Send Slack hook TODO(spicer): Add more information like email.
	slack.Notify("#events", fmt.Sprintf("(%s) New Snap!Clerk submission. Account: %d, Id: %d", os.Getenv("APP_ENV"), sc.AccountId, sc.Id))

	// Some logging
	services.Info(errors.New(fmt.Sprintf("New Snap!Clerk received. Account: %d, Id: %d", sc.AccountId, sc.Id)))

	return nil
}

//
// GetSnapClerkByAccountAndId by account and id.
//
func (db *DB) GetSnapClerkByAccountAndId(accountId uint, id uint) (SnapClerk, error) {
	// SnapClerk to return
	c := SnapClerk{}

	// Make query
	if db.New().Preload("File").Where("SnapClerkAccountId = ? AND SnapClerkId = ?", accountId, id).First(&c).RecordNotFound() {
		return SnapClerk{}, errors.New("SnapClerk entry not found.")
	}

	// Loop through and add the signed URLs to the files
	if len(c.File.Path) > 0 {
		c.File.Url = db.GetSignedFileUrl(c.File.Path)
	}

	if len(c.File.ThumbPath) > 0 {
		c.File.Thumb600By600Url = db.GetSignedFileUrl(c.File.ThumbPath)
	}

	// Return result
	return c, nil
}

//
// SnapClerkMonthlyUsage - Returns how many Snapclerks we have used this month.
//
func (db *DB) SnapClerkMonthlyUsage(accountId uint) int {
	// Return struct
	type r struct {
		Count int `json:"count"`
	}

	rt := r{}

	// Set the start date
	now := time.Now()
	start := now.Format("2006") + "-" + now.Format("01") + "-01"

	// SQL String
	sql := "SELECT COUNT(SnapClerkId) AS count FROM SnapClerk WHERE SnapClerkCreatedAt > ? AND SnapClerkAccountId = ?"

	// Run query
	db.New().Raw(sql, start, accountId).Scan(&rt)

	// Return happy.
	return rt.Count
}

//
// ConvertSnapclerkToLedger takes a snapclerk model saves a ledger entry.
//
func (db *DB) ConvertSnapclerkToLedger(sc SnapClerk) (Ledger, error) {
	ledger := Ledger{}

	// Create category
	ledger.Category = Category{Name: sc.Category, Type: "1", AccountId: sc.AccountId}

	// Create contact
	ledger.Contact = Contact{Name: sc.Contact, AccountId: sc.AccountId}

	// Build lat / Lon
	var lat float64
	var lon float64

	if l1, err := strconv.ParseFloat(sc.Lat, 64); err == nil {
		lat = l1
	}

	if l2, err := strconv.ParseFloat(sc.Lon, 64); err == nil {
		lon = l2
	}

	// Deal with labels.
	lbArray := []Label{}
	lbArray = append(lbArray, Label{Name: "Snap!Clerk", AccountId: sc.AccountId})
	lbs := strings.Split(sc.Labels, ",")

	for _, row := range lbs {
		st := strings.Trim(row, " ")

		if len(st) == 0 {
			continue
		}

		lb := Label{Name: st, AccountId: sc.AccountId}
		lbArray = append(lbArray, lb)
	}

	// Build the ledger.
	ledger.Date = sc.CreatedAt
	ledger.Files = append(ledger.Files, sc.File)
	ledger.AccountId = sc.AccountId
	ledger.AddedById = sc.AddedById
	ledger.Lat = lat
	ledger.Lon = lon
	ledger.Amount = sc.Amount
	ledger.Note = sc.Note
	ledger.Labels = lbArray

	// Save ledger entry.
	err := db.LedgerCreate(&ledger)

	if err != nil {
		return ledger, err
	}

	// Get the contact name.
	contactName := ledger.Contact.Name

	if len(contactName) == 0 {
		contactName = ledger.Contact.FirstName + " " + ledger.Contact.LastName
	}

	// Add to the activity log
	db.New().Create(&Activity{
		AccountId: ledger.AccountId,
		UserId:    ledger.AddedById,
		Action:    "expense",
		SubAction: "create",
		Name:      contactName,
		Amount:    ledger.Amount,
		LedgerId:  ledger.Id,
	})

	return ledger, nil
}

/* End File */
