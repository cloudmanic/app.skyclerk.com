//
// Date: 2019-09-16
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
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

// //
// // GetAccountById - Get a account by Id.
// //
// func (t *DB) GetAccountById(id uint) (Account, error) {
// 	var u Account
//
// 	if t.Where("id = ?", id).First(&u).RecordNotFound() {
// 		return u, errors.New("Record not found")
// 	}
//
// 	// Return the user.
// 	return u, nil
// }

/* End File */
