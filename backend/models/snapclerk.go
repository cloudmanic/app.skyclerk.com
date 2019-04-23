//
// Date: 2019-04-22
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"time"
)

// SnapClerk struct
type SnapClerk struct {
	Id           uint      `gorm:"primary_key;column:SnapClerkId" json:"id"`
	AccountId    uint      `gorm:"column:SnapClerkAccountId" sql:"not null" json:"account_id"`
	AddedById    uint      `gorm:"column:SnapClerkAddedById" sql:"not null" json:"_"`
	ReviewedById uint      `gorm:"column:SnapClerkReviewedById" sql:"not null" json:"_"`
	Status       string    `gorm:"column:SnapClerkStatus" sql:"not null;type:ENUM('Pending','Processed','Rejected');default:'Pending'" json:"status"`
	FileId       uint      `gorm:"column:SnapClerkFileId" sql:"not null" json:"_"`
	File         File      `gorm:"foreignkey:SnapClerkFileId" json:"file"`
	LedgerId     uint      `gorm:"column:SnapClerkLedgerId;index:SnapClerkLedgerId" sql:"not null" json:"ledger_id"`
	Amount       float64   `gorm:"column:SnapClerkAmount" sql:"not null;type:DECIMAL(12,2)" json:"amount"`
	Contact      string    `gorm:"column:SnapClerkContact" sql:"not null" json:"contact"`
	Category     string    `gorm:"column:SnapClerkCategory" sql:"not null" json:"category"`
	Labels       string    `gorm:"column:SnapClerkLabels" sql:"not null" json:"labels"`
	Note         string    `gorm:"column:SnapClerkNote" sql:"not null;type:TEXT" json:"note"`
	Lat          string    `gorm:"column:SnapClerkLat" sql:"not null;type:TEXT" json:"lat"`
	Lon          string    `gorm:"column:SnapClerkLon" sql:"not null;type:TEXT" json:"lon"`
	UpdatedAt    time.Time `gorm:"column:SnapClerkUpdatedAt" sql:"not null" json:"_"`
	CreatedAt    time.Time `gorm:"column:SnapClerkCreatedAt" sql:"not null" json:"created_at"`
	ProcessedAt  time.Time `gorm:"column:SnapClerkProcessedAt" sql:"not null" json:"processed_at"` // This has to be at the end GORM (auto update stuff)
}

//
// Set the table name.
//
func (SnapClerk) TableName() string {
	return "SnapClerk"
}

//
// SnapClerkCreate - Create a new SnapClerk entry.
//
func (db *DB) SnapClerkCreate(sc *SnapClerk) error {
	// Store this entry.
	db.Create(&sc)

	// TODO(spicer): Add to AppLog

	// TODO(spicer): Send email to customer

	// TODO(spicer): Send Slack hook

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

/* End File */
