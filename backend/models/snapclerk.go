package models

import (
	"time"
)

// SnapClerk struct
type SnapClerk struct {
	Id           uint      `gorm:"primary_key;column:SnapClerkId" json:"id"`
	AccountId    uint      `gorm:"column:SnapClerkAccountId" sql:"not null" json:"account_id"`
	AddedById    uint      `gorm:"column:SnapClerkAddedById" sql:"not null" json:"_"`
	ReviewedById uint      `gorm:"column:SnapClerkReviewedById" sql:"not null" json:"_"`
	Status       uint      `gorm:"column:SnapClerkStatus" sql:"not null;type:ENUM('Pending','Processed','Rejected');default:'Pending'" json:"status"`
	FileId       uint      `gorm:"column:SnapClerkFileId" sql:"not null" json:"_"`
	File         File      `gorm:"foreignkey:SnapClerkFileId" json:"File"`
	LedgerId     uint      `gorm:"column:SnapClerkLedgerId" sql:"not null" json:"_"`
	Ledger       File      `gorm:"foreignkey:SnapClerkLedgerId" json:"ledger"`
	Amount       float64   `gorm:"column:SnapClerkAmount" sql:"not null;type:DECIMAL(12,2)" json:"amount"`
	Contact      string    `gorm:"column:SnapClerkContact" sql:"not null" json:"contact"`
	Category     string    `gorm:"column:SnapClerkCategory" sql:"not null" json:"category"`
	Labels       string    `gorm:"column:SnapClerkLabels" sql:"not null;type:TEXT" json:"labels"`
	Note         string    `gorm:"column:SnapClerkNote" sql:"not null;type:TEXT" json:"note"`
	Lat          string    `gorm:"column:SnapClerkLat" sql:"not null;type:TEXT" json:"lat"`
	Lon          string    `gorm:"column:SnapClerkLon" sql:"not null;type:TEXT" json:"lon"`
	Paid         bool      `gorm:"column:SnapClerkPaid" sql:"not null" json:"paid"`
	UpdatedAt    time.Time `gorm:"column:SnapClerkUpdatedAt" sql:"not null" json:"updated_at"`
	CreatedAt    time.Time `gorm:"column:SnapClerkCreatedAt" sql:"not null" json:"created_at"`
	ProcessedAt  time.Time `gorm:"column:SnapClerkProcessedAt" sql:"not null" json:"processed_at"` // This has to be at the end GORM (auto update stuff)
}

//
// Set the table name.
//
func (SnapClerk) TableName() string {
	return "SnapClerk"
}

/* End File */
