//
// Date: 2018-03-20
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-21
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import "time"

// Account struct
type Account struct {
	AccountsId           uint      `gorm:"primary_key;column:AccountsId" json:"id"`
	AccountsAppId        uint      `gorm:"column:AccountsAppId" json:"_"`
	AccountsOwnerId      uint      `gorm:"column:AccountsOwnerId" json:"owner_id"`
	AccountsDisplayName  string    `gorm:"column:AccountsDisplayName" json:"display_name"`
	AccountsPlanId       uint      `gorm:"column:AccountsPlanId" json:"_"`
	AccountsAddress      string    `gorm:"column:AccountsAddress" sql:"not null;type:TEXT" json:"_"`
	AccountsCity         string    `gorm:"column:AccountsCity" json:"_"`
	AccountsState        string    `gorm:"column:AccountsState" json:"_"`
	AccountsZip          string    `gorm:"column:AccountsZip" json:"_"`
	AccountsCountry      string    `gorm:"column:AccountsCountry" json:"_"`
	AccountsLastActivity time.Time `gorm:"column:AccountsLastActivity" json:"_"`
	AccountsStripeId     string    `gorm:"column:AccountsStripeId" json:"_"`
	AccountsCardType     string    `gorm:"column:AccountsCardType" json:"_"`
	AccountsCardLast4    string    `gorm:"column:AccountsCardLast4" json:"_"`
	AccountsCardExpMonth string    `gorm:"column:AccountsCardExpMonth" json:"_"`
	AccountsCardExpYear  string    `gorm:"column:AccountsCardExpYear" json:"_"`
	AccountsSignupIp     string    `gorm:"column:AccountsSignupIp" json:"_"`
	AccountsUpdatedAt    time.Time `gorm:"column:AccountsUpdatedAt" sql:"not null" json:"_"`
	AccountsCreatedAt    time.Time `gorm:"column:AccountsCreatedAt" sql:"not null" json:"_"`
}

//
// Set the table name.
//
func (Account) TableName() string {
	return "Accounts"
}

/* End File */
