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

	// Account
	ClearAccount(accountId uint)
	DeleteAccount(accountId uint)
	GetAccountById(id uint) (Account, error)
	ValidateOwnerId(acct Account, accountId uint, objId uint, action string) error

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
	LoadDefaultCategories(accountId uint)
	DeleteCategoryByAccountAndId(accountId uint, categoryId uint) error
	GetCategoryByNameAndTypeAndAccountID(accountID uint, name string, catType string) (Category, error)
	ValidateDuplicateCategoryName(cat Category, accountId uint, objId uint, action string) error
	GetCategoryByAccountAndId(accountId uint, categoryId uint) (Category, error)
	GetCategoryUsageByAccount(accountId uint) []CategoryUsage

	// Contact
	GenerateAvatarsForAllMissing() error
	CreateContact(contact *Contact) error
	ConfirmContactAvatar(contact *Contact) error
	DeleteContactByAccountAndId(accountId uint, contactId uint) error
	GetContactByAccountAndId(accountId uint, conId uint) (Contact, error)
	ValidateContactNameOrFirstLast(contact Contact, accountId uint, objId uint, action string) error
	GenerateAvatarsForAllMissingWoker(jobs <-chan generateAvatarsWorkerJob, results chan<- int)

	// Label
	GetLabelByAccountAndId(accountId uint, labelId uint) (Label, error)
	DeleteLabelByAccountAndId(accountId uint, labelId uint) error
	ValidateDuplicateLabelName(obj Label, accountId uint, objId uint, action string) error
	GetLabelUsageByAccount(accountId uint) []LabelUsage
	GetLabelByAccountAndName(accountId uint, name string) (Label, error)

	// User
	GetUserById(id uint) (User, error)
	GetUserByEmail(email string) (User, error)
	ValidatePassword(password string) error
	ValidateEmailAddress(email string) error
	ResetUserPassword(id uint, password string) error
	ValidateUserLogin(email string, password string) error
	ValidateCreateUser(first string, last string, email string, googleAuth bool) error
	LoginUserByEmailPass(email string, password string, appId uint, userAgent string, ipAddress string) (User, Session, error)
	CreateUser(first string, last string, email string, password string, appId uint, userAgent string, ipAddress string) (User, error)

	// AcctToUsers
	GetUsersByAccount(accountId uint) []User

	// Forget password stuff
	GetUserFromToken(token string) (User, error)
	DeleteForgotPasswordByToken(token string) error
	DoResetPassword(user_email string, ip string) error

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
	SnapClerkMonthlyUsage(accountId uint) int
	ConvertSnapclerkToLedger(sc SnapClerk) (Ledger, error)
	GetSnapClerkByAccountAndId(accountId uint, id uint) (SnapClerk, error)

	// Billing
	GetBillingByAccountId(id uint) (Billing, error)

	// ConnectedAccounts
	GetConnectedAccountsByAccountIDAndConnection(accountID uint, connection string) (ConnectedAccounts, error)
}

/* End File */
