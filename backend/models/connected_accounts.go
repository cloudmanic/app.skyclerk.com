//
// Date: 2020-06-07
// Author: Spicer Matthews (spicer@skyclerk.com)
// Last Modified by: Spicer Matthews
// Copyright: 2020 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"errors"
	"time"
)

// ConnectedAccounts struct
type ConnectedAccounts struct {
	ID                      uint      `gorm:"primary_key" json:"id"`
	CreatedAt               time.Time `sql:"not null" json:"-"`
	UpdatedAt               time.Time `sql:"not null" json:"-"`
	AccountID               uint      `sql:"not null" json:"account_id"`
	Name                    string    `sql:"not null" json:"name"`
	Connection              string    `sql:"not null;type:ENUM('Stripe', 'Harvest');default:'Stripe'" json:"connection"`
	StripeUserID            string    `sql:"not null" json:"-"`
	StripeIncomeCategoryID  uint      `sql:"not null" json:"stripe_income_category_id"`
	StripeExpenseCategoryID uint      `sql:"not null" json:"stripe_expense_category_id"`
	StripePublishableKey    string    `sql:"not null" json:"stripe_publishable_key"`
	StripeScope             string    `sql:"not null" json:"stripe_scope"`
	StripeLastItem          int64     `sql:"not null" json:"stripe_last_item"`
	StripeLastSync          time.Time `sql:"not null" json:"stripe_last_sync"`
}

//
// GetConnectedAccountsByAccountIDAndConnection - get entry by account and connection.
//
func (t *DB) GetConnectedAccountsByAccountIDAndConnection(accountID uint, connection string) (ConnectedAccounts, error) {
	b := ConnectedAccounts{}

	// Look up the billing.
	if t.Where("account_id = ? AND connection = ?", accountID, connection).First(&b).RecordNotFound() {
		return b, errors.New("Record not found")
	}

	// Return the object.
	return b, nil
}

/* End File */
