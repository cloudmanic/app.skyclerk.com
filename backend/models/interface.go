//
// Date: 2018-03-20
// Author: spicer (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2019-01-13
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import "github.com/jinzhu/gorm"

// Datastore interface
type Datastore interface {
	// Gorm Functions
	New() *gorm.DB

	// Generic database functions
	Query(model interface{}, params QueryParam) error
	QueryMeta(model interface{}, params QueryParam) (QueryMetaData, error)
	QueryWithNoFilterCount(model interface{}, params QueryParam) (int, error)
	GetQueryMetaData(noLimitCount int, params QueryParam) QueryMetaData

	// Application
	ValidateClientIdGrantType(clientId string, grantType string) (Application, error)

	// Session
	GetByAccessToken(accessToken string) (Session, error)

	// Ledger
	LedgerCreate(ledger *Ledger) error
	LedgerUpdate(ledger *Ledger) error
	DeleteLedgerByAccountAndId(accountId uint, id uint) error
	GetLedgerByAccountAndId(accountId uint, id uint) (Ledger, error)
	AddFileToLedgerEntry(accountId uint, ledgerId uint, fileId uint) error
	ValidateLedgerContact(ledger Ledger, accountId uint, objId uint, action string) error
	ValidateLedgerCategory(ledger Ledger, accountId uint, objId uint, action string) error

	// Category
	DeleteCategoryByAccountAndId(accountId uint, categoryId uint) error
	ValidateDuplicateCategoryName(cat Category, accountId uint, objId uint, action string) error
	GetCategoryByAccountAndId(accountId uint, categoryId uint) (Category, error)

	// Contact
	DeleteContactByAccountAndId(accountId uint, contactId uint) error
	GetContactByAccountAndId(accountId uint, conId uint) (Contact, error)
	ValidateContactNameOrFirstLast(contact Contact, accountId uint, objId uint, action string) error

	// Label
	GetLabelByAccountAndId(accountId uint, labelId uint) (Label, error)
	DeleteLabelByAccountAndId(accountId uint, labelId uint) error
	ValidateDuplicateLabelName(obj Label, accountId uint, objId uint, action string) error

	// User
	GetUserById(id uint) (User, error)
	GetUserByEmail(email string) (User, error)
	ValidatePassword(password string) error
	ValidateEmailAddress(email string) error
	ValidateUserLogin(email string, password string) error
	LoginUserByEmailPass(email string, password string, appId uint, userAgent string, ipAddress string) (User, Session, error)

	// File
	GetSignedFileUrl(path string) string
	CleanFileName(fileName string) string
	StoreFile(accountId uint, filePath string) (File, error)
	GetFileByAccountAndId(accountId uint, id uint) (File, error)
	GetImageThumbNail(file *File, filePath string, width int, height int, cleanedFileName string) (string, error)
	GetPdfThumbNail(file *File, width int, height int, cleanedFileName string) (string, error)
	CreateAndStoreThumbnailImage(file *File, cleanedFileName string, filePath string, fileType string) error

	// SnapClerk
	SnapClerkCreate(sc *SnapClerk) error
}

/* End File */
