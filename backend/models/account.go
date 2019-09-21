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
