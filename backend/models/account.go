//
// Date: 2018-03-20
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-21
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Account struct
type Account struct {
	Id           uint      `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time `sql:"not null" json:"-"`
	UpdatedAt    time.Time `sql:"not null" json:"-"`
	OwnerId      uint      `sql:"not null" json:"owner_id"`
	Name         string    `sql:"not null" json:"name"`
	Address      string    `sql:"not null;type:TEXT" json:"-"`
	City         string    `sql:"not null" json:"-"`
	State        string    `sql:"not null" json:"-"`
	Zip          string    `sql:"not null" json:"-"`
	Country      string    `sql:"not null" json:"-"`
	Locale       string    `sql:"not null;default:'en-US'" json:"locale"` // BCP 47 language tag
	Currency     string    `sql:"not null;default:'USD'" json:"currency"` // The ISO 4217 currency code, such as USD for the US dollar and EUR for the euro.
	LastActivity time.Time `sql:"not null" json:"-"`
}

//
// Validate for this model.
//
func (a Account) Validate(db Datastore, action string, userId uint, accountId uint, objId uint) error {
	return validation.ValidateStruct(&a,
		// Name
		validation.Field(&a.Name,
			validation.Required.Error("The name field is required."),
			validation.Length(1, 50).Error("The name is too long."),
		),

		// Locale
		validation.Field(&a.Locale,
			validation.Required.Error("The locale field is required."),
		),

		// Currency
		validation.Field(&a.Currency,
			validation.Required.Error("The currency field is required."),
		),

		// OwnerId
		validation.Field(&a.OwnerId,
			validation.Required.Error("The owner_id field is required."),
			validation.By(func(value interface{}) error { return db.ValidateOwnerId(a, accountId, objId, action) }),
		),
	)
}

//
// ValidateOwnerId - is this a valid owner id?
//
func (db *DB) ValidateOwnerId(acct Account, accountId uint, objId uint, action string) error {
	// Default error message
	const errMsg = "Invalid owner_id was posted."

	// Make sure this user is part of the account.
	c := AcctToUsers{}
	if db.New().Where("acct_id = ? AND user_id = ?", accountId, acct.OwnerId).First(&c).RecordNotFound() {
		return errors.New(errMsg)
	}

	// All good in the hood
	return nil
}

//
// GetAccountById - Get a account by Id.
//
func (t *DB) GetAccountById(id uint) (Account, error) {
	var u Account

	if t.Where("id = ?", id).First(&u).RecordNotFound() {
		return u, errors.New("Record not found")
	}

	// Return the user.
	return u, nil
}

/* End File */
