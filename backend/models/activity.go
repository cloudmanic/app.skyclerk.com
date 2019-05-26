//
// Date: 2018-05-26
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

// Models
package models

import (
	"time"
)

// Activity struct
type Activity struct {
	Id          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	AccountId   uint    `sql:"not null;index:account_id" json:"account_id"`
	UserId      uint    `sql:"not null;index:user_id"`
	Action      string  `sql:"not null;type:ENUM('income', 'expense', 'contact', 'category', 'label', 'snapclerk', 'other');default:'other'" json:"action"`
	SubAction   string  `sql:"not null;type:ENUM('create', 'update', 'delete', 'other');default:'other'" json:"sub_action"`
	Name        string  `sql:"not null" json:"name"`
	Amount      float64 `sql:"not null;type:DECIMAL(12,2)" json:"amount"`
	LedgerId    uint    `sql:"not null;index:ledger_id" json:"ledger_id"`
	ContactId   uint    `sql:"not null;index:contact_id" json:"contact_id"`
	LabelId     uint    `sql:"not null;index:label_id" json:"label_id"`
	CategoryId  uint    `sql:"not null;index:category_id" json:"category_id"`
	SnapClerkId uint    `sql:"not null;index:snapclerk_id" json:"snapclerk_id"`
}

/* End File */
