//
// Date: 2018-03-20
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2019-01-13
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Contact struct {
	Id            uint      `gorm:"primary_key;column:ContactsId" json:"id"`
	AccountId     uint      `gorm:"column:ContactsAccountId" sql:"not null" json:"account_id"`
	UpdatedAt     time.Time `gorm:"column:ContactsUpdatedAt" sql:"not null" json:"created_at"`
	CreatedAt     time.Time `gorm:"column:ContactsCreatedAt" sql:"not null" json:"updated_at"`
	Name          string    `gorm:"column:ContactsName" sql:"not null" json:"name"`
	FirstName     string    `gorm:"column:ContactsFirstName" sql:"not null" json:"first_name"`
	LastName      string    `gorm:"column:ContactsLastName" sql:"not null" json:"last_name"`
	Address       string    `gorm:"column:ContactsAddress" sql:"not null" json:"address"`
	City          string    `gorm:"column:ContactsCity" sql:"not null" json:"city"`
	State         string    `gorm:"column:ContactsState" sql:"not null" json:"state"`
	Zip           string    `gorm:"column:ContactsZip" sql:"not null" json:"zip"`
	Phone         string    `gorm:"column:ContactsPhone" sql:"not null" json:"phone"`
	Fax           string    `gorm:"column:ContactsFax" sql:"not null" json:"fax"`
	Website       string    `gorm:"column:ContactsWebsite" sql:"not null" json:"website"`
	AccountNumber string    `gorm:"column:ContactsAccountNumber" sql:"not null" json:"account_number"`
	Avatar        string    `gorm:"column:ContactsAvatar" sql:"not null" json:"_"`
	AvatarChecked string    `gorm:"column:ContactsAvatarChecked" sql:"not null;type:ENUM('Yes', 'No');default:'No'" json:"avatar_checked"`
	Email         string    `gorm:"column:ContactsEmail" sql:"not null" json:"email"`
	Twitter       string    `gorm:"column:ContactsTwitter" sql:"not null" json:"twitter"`
	Facebook      string    `gorm:"column:ContactsFacebook" sql:"not null" json:"facebook"`
	Linkedin      string    `gorm:"column:ContactsLinkedin" sql:"not null" json:"linkedin"`
	Type          string    `gorm:"column:ContactsType" sql:"not null;type:ENUM('Customer', 'Vendor', 'Both');default:'Both'" json:"_"`
	HrId          uint64    `gorm:"column:ContactsHrId" sql:"not null" json:"_"`
	PricingPlanId uint      `gorm:"column:ContactsPricingPlanId" sql:"not null" json:"_"`
	GatewayId     uint      `gorm:"column:ContactsGatewayId" sql:"not null" json:"_"`
	CardMask      string    `gorm:"column:ContactsCardMask" sql:"not null" json:"_"`
	CardType      string    `gorm:"column:ContactsCardType" sql:"not null" json:"_"`
	CardExpire    string    `gorm:"column:ContactsCardExpire" sql:"not null" json:"_"`
	Country       string    `gorm:"column:ContactsCountry" sql:"not null" json:"country"`
}

//
// Set the table name.
//
func (Contact) TableName() string {
	return "Contacts"
}

//
// Validate for this model.
//
func (a Contact) Validate(db Datastore, action string, userId uint, accountId uint, objId uint) error {

	return validation.ValidateStruct(&a,

		validation.Field(&a.Name,
			validation.By(func(value interface{}) error { return db.ValidateContactNameOrFirstLast(a, accountId, objId, action) }),
		),

		// validation.Field(&a.Type,
		// 	validation.Required.Error("The type field is required."),
		// 	validation.In("1", "2").Error("The type field must be 1, or 2."),
		// ),
	)
}

//
// Validate Duplicate Contact Name, First, Last
//
func (db *DB) ValidateContactNameOrFirstLast(contact Contact, accountId uint, objId uint, action string) error {

	// trim any white space
	contactName := strings.Trim(contact.Name, " ")
	contactFirstName := strings.Trim(contact.FirstName, " ")
	contactLastName := strings.Trim(contact.LastName, " ")

	// Make sure this category is not already in use.
	if action == "create" {

		c := Contact{}

		// IF we have a contact name we validate against that.
		if len(contactName) > 0 {

			if !db.New().Where("ContactsAccountId = ? AND ContactsName = ?", accountId, contactName).First(&c).RecordNotFound() {
				return errors.New("Contact name is already in use.")
			}

			// Double check casing
			if strings.ToLower(contactName) == strings.ToLower(c.Name) {
				return errors.New("Contact name is already in use.")
			}

		} else if (len(contactFirstName) > 0) && (len(contactLastName) > 0) { // Validate first / last

			if !db.New().Where("ContactsAccountId = ? AND ContactsFirstName = ? AND ContactsLastName = ?", accountId, contactFirstName, contactLastName).First(&c).RecordNotFound() {
				return errors.New("Contact first and last name is already in use.")
			}

			// Double check casing
			if strings.ToLower(contactFirstName) == strings.ToLower(c.FirstName) {
				return errors.New("Contact first and last name is already in use.")
			}

			if strings.ToLower(contactLastName) == strings.ToLower(c.LastName) {
				return errors.New("Contact first and last name is already in use.")
			}
		}

	} else if action == "update" {

		c := Contact{}

		if !db.New().Where("ContactsAccountId = ? AND ContactsName = ? AND ContactsFirstName = ? AND ContactsLastName = ?", accountId, contactName, contactFirstName, contactLastName).First(&c).RecordNotFound() {

			// Make sure it is not the same id as the one we are updating
			if c.Id != objId {
				return errors.New("Contact company name, first, and last name is already in use.")
			}
		}

		// Double check casing
		if (c.Id != objId) && (strings.ToLower(contactName) == strings.ToLower(c.Name) && (strings.ToLower(contactFirstName) == strings.ToLower(c.FirstName)) && (strings.ToLower(contactLastName) == strings.ToLower(c.LastName))) {
			return errors.New("Contact company name, first, and last name is already in use.")
		}

	}

	// All good in the hood
	return nil
}

//
// GetContactByAccountAndId - Return a contact by account and id.
//
func (db *DB) GetContactByAccountAndId(accountId uint, conId uint) (Contact, error) {
	l := Contact{}

	// Make query
	if db.New().Where("ContactsAccountId = ? AND ContactsId = ?", accountId, conId).First(&l).RecordNotFound() {
		return Contact{}, errors.New("Contact not found.")
	}

	// Return result
	return l, nil
}

//
// DeleteContactByAccountAndId - Delete a contact by account and id.
//
func (db *DB) DeleteContactByAccountAndId(accountId uint, contactId uint) error {
	// Make query to see if we have ledger entries with this contact
	if !db.New().Where("LedgerAccountId = ? AND LedgerContactId = ?", accountId, contactId).First(&Ledger{}).RecordNotFound() {
		return errors.New("Can not delete contact. It is in use by a ledger entry.")
	}

	// Make query
	db.New().Where("ContactsAccountId = ? AND ContactsId = ?", accountId, contactId).Delete(Contact{})

	// Return result
	return nil
}

/* End File */
