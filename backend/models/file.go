package models

import (
	"time"

	"github.com/asaskevich/govalidator"
)

// File struct
type File struct {
	Id        uint      `gorm:"primary_key;column:FilesId" json:"id"`
	AccountId uint      `gorm:"column:FilesAccountId" sql:"not null" json:"account_id"`
	UpdatedAt time.Time `gorm:"column:FilesUpdatedAt" sql:"not null" json:"_"`
	CreatedAt time.Time `gorm:"column:FilesCreatedAt" sql:"not null" json:"_"`
	Host      string    `gorm:"column:FilesHost" sql:"not null" json:"_"`
	Name      string    `gorm:"column:FilesName" sql:"not null" json:"name"`
	Path      string    `gorm:"column:FilesPath" sql:"not null" json:"_"`
	ThumbPath string    `gorm:"column:FilesThumbPath" sql:"not null" json:"_"`
	Type      string    `gorm:"column:FilesType" sql:"not null" json:"type"`
	Hash      string    `gorm:"column:FilesHash" sql:"not null" json:"_"`
	Size      int64     `gorm:"column:FilesSize" sql:"not null" json:"size"`
}

//
// Set the table name.
//
func (File) TableName() string {
	return "Files"
}

//
// CleanFileName returns a cleaned-up filename that is safe to use.
//
func (t *DB) CleanFileName(fileName string) string {
	// This already checks for ";" so that ";;DROP TABLE x;" can't be part of the name.
	return govalidator.SafeFileName(fileName)
}

/* End File */
