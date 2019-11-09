//
// Date: 2018-03-20
// Author: Spicer Matthews (spicer@cloudmanic.com)
// Last Modified by: Spicer Matthews
// Last Modified: 2018-03-21
// Copyright: 2017 Cloudmanic Labs, LLC. All rights reserved.
//

package models

import (
	"time"
)

// AcctToBilling struct
type AcctToBilling struct {
	Id        uint      `gorm:"primary_key" json:"-"`
	CreatedAt time.Time `sql:"not null" json:"-"`
	UpdatedAt time.Time `sql:"not null" json:"-"`
	AccountId uint      `sql:"not null"  json:"-"`
	BillingId uint      `sql:"not null"  json:"-"`
}

/* End File */
