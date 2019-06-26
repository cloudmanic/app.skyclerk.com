//
// Date: 2018-05-26
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

// Models
package models

import (
	"fmt"
	"time"
)

// Activity struct
type Activity struct {
	Id          uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	AccountId   uint    `sql:"not null;index:account_id" json:"account_id"`
	UserId      uint    `sql:"not null;index:user_id"`
	User        User    `gorm:"foreignkey:UserId" json:"user"`
	Action      string  `sql:"not null;type:ENUM('income', 'expense', 'contact', 'category', 'label', 'snapclerk', 'other');default:'other'" json:"action"`
	SubAction   string  `sql:"not null;type:ENUM('create', 'update', 'delete', 'other');default:'other'" json:"sub_action"`
	Name        string  `sql:"not null" json:"name"`
	Amount      float64 `sql:"not null;type:DECIMAL(12,2)" json:"amount"`
	LedgerId    uint    `sql:"not null;index:ledger_id" json:"ledger_id"`
	ContactId   uint    `sql:"not null;index:contact_id" json:"contact_id"`
	LabelId     uint    `sql:"not null;index:label_id" json:"label_id"`
	CategoryId  uint    `sql:"not null;index:category_id" json:"category_id"`
	SnapClerkId uint    `sql:"not null;index:snapclerk_id" json:"snapclerk_id"`
	Message     string  `gorm:"-" json:"message"`
}

//
// SetMessage based on the fields we pass in.
//
func (a *Activity) SetMessage() {
	subAction := ""
	userName := "Unknown"

	// Make sure we have a user
	if a.User.Id > 0 {
		userName = a.User.FirstName
	}

	if a.SubAction == "create" {
		subAction = "created"
	}

	// See if this is a ledger activity. - Spicer, Added a ledger entry of -2325.20 for Bank of America.
	if a.LedgerId > 0 {
		a.Message = fmt.Sprintf("%s, %s an %s ledger entry of %.2f for %s.", userName, subAction, a.Action, a.Amount, a.Name)
	}
}

/* End File */
