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
	Id          uint      `gorm:"primary_key" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	AccountId   uint      `sql:"not null;index:account_id" json:"account_id"`
	UserId      uint      `sql:"not null;index:user_id" json:"user_id"`
	User        User      `json:"user"`
	Action      string    `sql:"not null;type:ENUM('income', 'expense', 'contact', 'category', 'label', 'snapclerk', 'other');default:'other'" json:"action"`
	SubAction   string    `sql:"not null;type:ENUM('create', 'update', 'delete', 'other');default:'other'" json:"sub_action"`
	Name        string    `sql:"not null" json:"name"`
	Amount      float64   `sql:"not null;type:DECIMAL(12,2)" json:"amount"`
	LedgerId    uint      `sql:"not null;index:ledger_id" json:"ledger_id"`
	Ledger      Ledger    `json:"ledger"`
	ContactId   uint      `sql:"not null;index:contact_id" json:"contact_id"`
	LabelId     uint      `sql:"not null;index:label_id" json:"label_id"`
	CategoryId  uint      `sql:"not null;index:category_id" json:"category_id"`
	SnapClerkId uint      `sql:"not null;index:snapclerk_id" json:"snapclerk_id"`
	Message     string    `gorm:"-" json:"message"`
}

//
// SetMessage based on the fields we pass in.
//
func (a *Activity) SetMessage() {
	subAction := ""
	mixWord := ""
	userName := "Unknown"
	contactName := ""

	// Get contact
	if a.LedgerId > 0 {
		contactName = a.Ledger.Contact.Name

		if len(contactName) == 0 {
			contactName = a.Ledger.Contact.FirstName + " " + a.Ledger.Contact.LastName
		}
	}

	// Make sure we have a user
	if a.User.Id > 0 {
		userName = a.User.FirstName
	}

	if a.SubAction == "create" {
		mixWord = "from"
		subAction = "created"
	}

	if a.SubAction == "update" {
		mixWord = "for"
		subAction = "updated"
	}

	if a.SubAction == "delete" {
		mixWord = "for"
		subAction = "deleted"
	}

	// See if this is a ledger activity. - Spicer, Added a ledger entry of -2325.20 for Bank of America.
	if a.LedgerId > 0 {
		a.Message = fmt.Sprintf("%s %s an %s ledger entry of %.2f %s %s.", userName, subAction, a.Action, a.Amount, mixWord, a.Name)
	}

	// See if this is a snapclerk activity. - Spicer, Uploaded a new receipt to be processed.
	if a.SnapClerkId > 0 {
		// Create
		if a.SubAction == "create" {
			a.Message = fmt.Sprintf("%s uploaded a new receipt to be processed.", userName)
		}

		// Update
		if a.SubAction == "update" {
			// Success
			if a.LedgerId > 0 {
				a.Message = fmt.Sprintf("%s's uploaded receipt has been processed for %s.", userName, contactName)
			}
		}
	}
}

/* End File */
