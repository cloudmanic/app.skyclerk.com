//
// Date: 2018-03-20
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-22
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
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
// Validate for this model.
//
func (a Ledger) Validate(db Datastore, action string, userId uint, accountId uint, objId uint) error {
	return validation.ValidateStruct(&a,

		validation.Field(&a.Amount,
			validation.Required.Error("The amount field is required."),
		),

		validation.Field(&a.Date,
			validation.Required.Error("The date field is required."),
		),

		validation.Field(&a.Category,
			validation.By(func(value interface{}) error { return db.ValidateLedgerCategory(a, accountId, objId, action) }),
		),

		validation.Field(&a.Contact,
			validation.By(func(value interface{}) error { return db.ValidateLedgerContact(a, accountId, objId, action) }),
		),
	)
}

//
// ValidateLedgerContact - Make sure all is good.
//
func (db *DB) ValidateLedgerContact(ledger Ledger, accountId uint, objId uint, action string) error {
	const errMsg1 = "Contact name is required."
	const errMsg2 = "Contact first and last name is required."

	if len(strings.Trim(ledger.Contact.Name, " ")) <= 0 {
		return errors.New(errMsg1)
	}

	if (len(strings.Trim(ledger.Contact.FirstName, " ")) <= 0) || (len(strings.Trim(ledger.Contact.LastName, " ")) <= 0) {
		return errors.New(errMsg2)
	}

	// All good in the hood
	return nil
}

//
// ValidateLedgerCategory - Make sure all is good.
//
func (db *DB) ValidateLedgerCategory(ledger Ledger, accountId uint, objId uint, action string) error {
	const errMsg1 = "Category name is required."
	const errMsg2 = "Category type is required."

	if len(strings.Trim(ledger.Category.Name, " ")) <= 0 {
		return errors.New(errMsg1)
	}

	if len(strings.Trim(ledger.Category.Type, " ")) <= 0 {
		return errors.New(errMsg2)
	}

	// All good in the hood
	return nil
}

//
// LedgerCreate - Create a new ledger entry.
//
func (db *DB) LedgerCreate(ledger *Ledger) error {
	// Make sure there is no funny biz with account ids. We make sure ledger.AccountId is always set correctly
	ledger.Contact.AccountId = ledger.AccountId
	ledger.Category.AccountId = ledger.AccountId

	// Trim Contact
	ledger.Contact.Name = strings.Trim(ledger.Contact.Name, " ")
	ledger.Contact.FirstName = strings.Trim(ledger.Contact.FirstName, " ")
	ledger.Contact.LastName = strings.Trim(ledger.Contact.LastName, " ")

	// Trim Category
	ledger.Category.Name = strings.Trim(ledger.Category.Name, " ")
	ledger.Category.Type = strings.Trim(ledger.Category.Type, " ")

	// Setup the contact. If we have a ledger.Contact.Id we assume we are not adding the contact on insert.
	if ledger.Contact.Id == 0 {
		if (len(ledger.Contact.FirstName) > 0) || (len(ledger.Contact.LastName) > 0) {
			db.Where("ContactsAccountId = ? AND ContactsName = ?", ledger.AccountId, ledger.Contact.Name).Or("ContactsAccountId = ? AND ContactsFirstName = ? AND ContactsLastName = ?", ledger.AccountId, ledger.Contact.FirstName, ledger.Contact.LastName).FirstOrCreate(&ledger.Contact)
		} else {
			db.Where("ContactsAccountId = ? AND ContactsName = ?", ledger.AccountId, ledger.Contact.Name).FirstOrCreate(&ledger.Contact)
		}
	}

	// Setup the category. Add the Id if we do not pass one in.
	if ledger.Category.Id == 0 {
		db.Where("CategoriesAccountId = ? AND CategoriesName = ? AND CategoriesType = ?", ledger.AccountId, ledger.Category.Name, ledger.Category.Type).FirstOrCreate(&ledger.Category)
	}

	// Setup the labels
	for key, row := range ledger.Labels {
		ledger.Labels[key].AccountId = ledger.AccountId
		db.Where("LabelsAccountId = ? AND LabelsName = ?", ledger.AccountId, strings.Trim(row.Name, " ")).FirstOrCreate(&ledger.Labels[key])
	}

	// Store this ledger entry.
	db.Create(&ledger)

	// Add additional data to lookups TODO: remove this once we retire PHP app.
	db.Model(LabelsToLedger{}).Where("LabelsToLedgerLedgerId = ?", ledger.Id).Updates(LabelsToLedger{LabelsToLedgerAccountId: ledger.AccountId, LabelsToLedgerCreatedAt: time.Now()})

	return nil
}

//
// GetLedgerByAccountAndId by account and id.
//
func (db *DB) GetLedgerByAccountAndId(accountId uint, id uint) (Ledger, error) {
	// Ledger to return
	c := Ledger{}

	// Make query
	if db.New().Preload("Contact").Preload("Category").Preload("Labels").Where("LedgerAccountId = ? AND LedgerId = ?", accountId, id).First(&c).RecordNotFound() {
		return Ledger{}, errors.New("Ledger entry not found.")
	}

	// Return result
	return c, nil
}

//
// DeleteLedgerByAccountAndId - Delete a label by account and id.
//
func (db *DB) DeleteLedgerByAccountAndId(accountId uint, id uint) error {
	// Make query to delete
	db.New().Where("LedgerAccountId = ? AND LedgerId = ?", accountId, id).Delete(Ledger{})

	// Delete from look up table.
	db.New().Where("LabelsToLedgerAccountId = ? AND LabelsToLedgerLedgerId = ?", accountId, id).Delete(LabelsToLedger{})

	// Return result
	return nil
}

/* End File */
