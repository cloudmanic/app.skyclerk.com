//
// Date: 2019-09-14
// Author: Spicer Matthews (spicer@skyclerk.com)
// Copyright: 2019 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"
)

// Invite struct
type Invite struct {
	Id        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `sql:"not null" json:"created_at"`
	UpdatedAt time.Time `sql:"not null" json:"updated_at"`
	AccountId uint      `sql:"not null;index:AccountId" json:"account_id"`
	Email     string    `sql:"not null" json:"email"`
	FirstName string    `sql:"not null" json:"first_name"`
	LastName  string    `sql:"not null" json:"last_name"`
	Message   string    `sql:"not null;type:TEXT" json:"message"`
	Token     string    `sql:"not null" json:"-"`
	ExpiresAt time.Time `sql:"not null" json:"expires_at"`
}

/* End File */
