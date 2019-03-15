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
	Count(model interface{}, params QueryParam) (uint, error)
	Query(model interface{}, params QueryParam) error
	QueryMeta(model interface{}, params QueryParam) (QueryMetaData, error)
	QueryWithNoFilterCount(model interface{}, params QueryParam) (int, error)
	GetQueryMetaData(limitCount int, noLimitCount int, params QueryParam) QueryMetaData

	// Ledger
	LedgerCreate(ledger *Ledger) error

	// Category
	DeleteCategoryByAccountAndId(accountId uint, categoryId uint) error
	ValidateDuplicateCategoryName(cat Category, accountId uint, objId uint, action string) error
	GetCategoryByAccountAndId(accountId uint, categoryId uint) (Category, error)

	// Contact
	DeleteContactByAccountAndId(accountId uint, contactId uint) error
	GetContactByAccountAndId(accountId uint, conId uint) (Contact, error)
	ValidateContactNameOrFirstLast(contact Contact, accountId uint, objId uint, action string) error

	// Labels
	GetLabelByAccountAndId(accountId uint, labelId uint) (Label, error)
	DeleteLabelByAccountAndId(accountId uint, labelId uint) error
	ValidateDuplicateLabelName(obj Label, accountId uint, objId uint, action string) error
}

/* End File */
