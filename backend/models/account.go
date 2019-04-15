//
// Date: 2018-03-20
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-21
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import "time"

// Account struct
type Account struct {
	Id           uint      `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
	OwnerId      uint      `sql:"not null" json:"owner_id"`
	Name         string    `sql:"not null" json:"name"`
	PlanId       uint      `sql:"not null" json:"_"`
	Address      string    `sql:"not null;type:TEXT" json:"_"`
	City         string    `sql:"not null" json:"_"`
	State        string    `sql:"not null" json:"_"`
	Zip          string    `sql:"not null" json:"_"`
	Country      string    `sql:"not null" json:"_"`
	LastActivity time.Time `sql:"not null" json:"_"`
	StripeId     string    `sql:"not null" json:"_"`
	CardType     string    `sql:"not null" json:"_"`
	CardLast4    string    `sql:"not null" json:"_"`
	CardExpMonth string    `sql:"not null" json:"_"`
	CardExpYear  string    `sql:"not null" json:"_"`
	SignupIp     string    `sql:"not null" json:"_"`
}

/* End File */
