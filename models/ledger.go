//
// Date: 2018-03-20
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-22
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"
)

type Ledger struct {
	Id               uint      `gorm:"primary_key;column:LedgerId" json:"id"`
	AccountId        uint      `gorm:"column:LedgerAccountId" sql:"not null" json:"account_id"`
	UpdatedAt        time.Time `gorm:"column:LedgerUpdatedAt" sql:"not null" json:"_"`
	CreatedAt        time.Time `gorm:"column:LedgerCreatedAt" sql:"not null" json:"_"`
	ContactId        uint      `gorm:"column:LedgerContactId" sql:"not null" json:"contact_id"`
	Contact          Contact   `gorm:"foreignkey:LedgerContactId" json:"contact"`
	Date             time.Time `gorm:"column:LedgerDate" sql:"not null" json:"date"`
	AddedById        uint      `gorm:"column:LedgerAddedById" sql:"not null" json:"added_by_id"`
	Amount           float64   `gorm:"column:LedgerAmount" sql:"type:DECIMAL(12,2)" json:"amount"`
	CategoryId       uint      `gorm:"column:LedgerCategoryId" sql:"not null" json:"category_id"`
	Category         Category  `gorm:"foreignkey:LedgerCategoryId" json:"category"`
	Note             string    `gorm:"column:LedgerNote" sql:"not null;type:TEXT" json:"note"`
	ShoeboxedId      string    `gorm:"column:LedgerShoeboxedId" sql:"not null" json:"_"`
	ShoeboxedImage   string    `gorm:"column:LedgerShoeboxedImage" sql:"not null" json:"_"`
	FreshBooksId     string    `gorm:"column:LedgerFreshBooksId" sql:"not null" json:"_"`
	AirBnbHash       string    `gorm:"column:LedgerAirBnbHash" sql:"not null" json:"_"`
	AuthGatewayToken string    `gorm:"column:LedgerAuthGatewayToken" sql:"not null" json:"_"`
	StripeId         string    `gorm:"column:LedgerStripeId" sql:"not null" json:"_"`
	Labels           []Label   `gorm:"many2many:LabelsToLedger;association_foreignkey:LabelsId;foreignkey:LedgerId;association_jointable_foreignkey:LabelsToLedgerLabelId;jointable_foreignkey:LabelsToLedgerLedgerId" sql:"not null" json:"labels"`
}

//
// Set the table name.
//
func (Ledger) TableName() string {
	return "Ledger"
}

//
// Create a new ledger entry.
//
func (db *DB) LedgerCreate(ledger *Ledger) error {

	// Setup the contact. If we have a ledger.Contact.Id we assume we are not adding the contact on insert.
	if ledger.Contact.Id == 0 {
		if (len(ledger.Contact.FirstName) > 0) || (len(ledger.Contact.LastName) > 0) {
			db.Where("ContactsAccountId = ? AND ContactsName = ?", ledger.Contact.AccountId, ledger.Contact.Name).Or("ContactsAccountId = ? AND ContactsFirstName = ? AND ContactsLastName = ?", ledger.Contact.AccountId, ledger.Contact.FirstName, ledger.Contact.LastName).FirstOrCreate(&ledger.Contact)
		} else {
			db.Where("ContactsAccountId = ? AND ContactsName = ?", ledger.Contact.AccountId, ledger.Contact.Name).FirstOrCreate(&ledger.Contact)
		}
	}

	// Setup the category
	if ledger.Category.Id == 0 {
		db.Where("CategoriesAccountId = ? AND CategoriesName = ? AND CategoriesType = ?", ledger.Category.AccountId, ledger.Category.Name, ledger.Category.Type).FirstOrCreate(&ledger.Category)
	}

	// Setup the labels
	for key, row := range ledger.Labels {
		db.Where("LabelsAccountId = ? AND LabelsName = ?", row.AccountId, row.Name).FirstOrCreate(&ledger.Labels[key])
	}

	// Store this ledger entry.
	db.Create(&ledger)

	// Add additional data to lookups TODO: remove this once we retire PHP app.
	db.Model(LabelsToLedger{}).Where("LabelsToLedgerLedgerId = ?", ledger.Id).Updates(LabelsToLedger{LabelsToLedgerAccountId: ledger.AccountId, LabelsToLedgerCreatedAt: time.Now()})

	return nil
}

/* End File */
