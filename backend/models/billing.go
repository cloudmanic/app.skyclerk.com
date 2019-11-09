//
// Date: 2019-09-16
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"time"
)

// Billing struct
type Billing struct {
	Id                 uint      `gorm:"primary_key" json:"id"`
	CreatedAt          time.Time `sql:"not null" json:"created_at"`
	UpdatedAt          time.Time `sql:"not null" json:"updated_at"`
	StripeCustomer     string    `sql:"not null" json:"-"`
	StripeSubscription string    `sql:"not null" json:"-"`
	Status             string    `sql:"not null;type:ENUM('Active', 'Disable', 'Delinquent', 'Expired', 'Trial');default:'Trial'" json:"status"`
	TrialExpire        time.Time `sql:"not null" json:"trial_expire"`
}

//
// GetBillingByAccountId - Get a billing by account id
//
func (t *DB) GetBillingByAccountId(id uint) (Billing, error) {
	var b Billing
	account := Account{}

	// Find in loop up table.
	if t.Where("id = ?", id).First(&account).RecordNotFound() {
		return b, errors.New("Account not found")
	}

	// Look up the billing.
	if t.Where("id = ?", account.BillingId).First(&b).RecordNotFound() {
		return b, errors.New("Record not found")
	}

	// Return the user.
	return b, nil
}

/* End File */
